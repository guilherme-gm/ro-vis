package skill

import (
	"testing"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
	"github.com/guilherme-gm/ro-vis/extractor/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestSkillLoader_loadJobs_InsertUpdateDelete(t *testing.T) {
	update := domain.Update{
		Date: time.Date(2012, time.January, 2, 0, 0, 0, 0, time.UTC),
		Changes: []domain.UpdateChange{
			{Patch: "testfiles", File: "data/luafiles514/lua files/skillinfoz/JobInheritList.lub"},
		},
	}

	// Seed existing jobs to trigger:
	// - Update: existing JT_SWORDMAN with wrong JobId and FirstUpdate set (should keep FirstUpdate)
	// - Delete: a job not present in file (TO_DELETE)
	// - Insert: JT_NOVICE not present in existing, present in file
	existing := []domain.SkillJob{
		{Constant: "JT_SWORDMAN", JobId: 999, FirstUpdate: "2012-01-01", LastUpdate: "2012-01-01", Deleted: false},
		{Constant: "TO_DELETE", JobId: 12345, FirstUpdate: "2012-01-01", LastUpdate: "2012-01-01", Deleted: false},
	}
	updater := loaders.NewUpdater[string](existing)

	loader := &SkillLoader{}
	loader.loadJobs(testfiles.Root, update, updater)

	// Verify inserts (JT_NOVICE should be inserted with values from file and timestamps set to update.Name())
	insertHasNovice := false
	for _, it := range updater.ValuesToInsert {
		if it.Constant == "JT_NOVICE" {
			insertHasNovice = true
			assert.Equal(t, int32(0), it.JobId)
			assert.Equal(t, update.Name(), it.FirstUpdate)
			assert.Equal(t, update.Name(), it.LastUpdate)
			assert.False(t, it.Deleted)
			break
		}
	}
	assert.True(t, insertHasNovice, "expected JT_NOVICE to be inserted")

	// Verify updates include JT_SWORDMAN corrected from 999 -> 1 and timestamps adjusted
	updateHasSwordman := false
	for _, it := range updater.ValuesToUpdate {
		if it.Constant == "JT_SWORDMAN" {
			updateHasSwordman = true
			assert.Equal(t, int32(1), it.JobId)
			// FirstUpdate should be preserved since it already existed
			assert.Equal(t, "2012-01-01", it.FirstUpdate)
			assert.Equal(t, update.Name(), it.LastUpdate)
			assert.False(t, it.Deleted)
			break
		}
	}
	assert.True(t, updateHasSwordman, "expected JT_SWORDMAN to be updated")

	// Verify delete: existing job not present in file should be marked deleted
	deleted, ok := updater.DirtyValues["TO_DELETE"]
	if assert.True(t, ok, "expected TO_DELETE to be edited for deletion") {
		assert.True(t, deleted.Deleted)
		assert.Equal(t, update.Name(), deleted.LastUpdate)
	}
}

func TestSkillLoader_loadSkillIDs_InsertUpdateDelete(t *testing.T) {
	update := domain.Update{
		Date: time.Date(2012, time.January, 2, 0, 0, 0, 0, time.UTC),
		Changes: []domain.UpdateChange{
			{Patch: "testfiles", File: "data/luafiles514/lua files/skillinfoz/SkillID.lub"},
		},
	}

	existing := []domain.Skill{
		{SkillID: 1, Constant: domain.NewNullableString("WRONG_CONST"), Deleted: false},
		{SkillID: 45, Constant: domain.NewNullableString("AC_CONCENTRATION"), Deleted: false},
		{SkillID: 9999, Constant: domain.NewNullableString("TO_DELETE"), Deleted: false},
	}
	skillUpdater := loaders.NewUpdater(existing)

	loader := &SkillLoader{}
	loader.loadSkillIDs(testfiles.Root, update, skillUpdater)

	// Verify unchanged
	_, exists := skillUpdater.DirtyValues[45]
	assert.False(t, exists, "expected AC_CONCENTRATION (45) to be unchanged")

	// Verify insertion
	insertHas := false
	for _, it := range skillUpdater.ValuesToInsert {
		if it.SkillID == 24 {
			insertHas = true
			assert.Equal(t, "AL_RUWACH", it.Constant.String)
			assert.True(t, it.Constant.Valid)
			assert.False(t, it.Deleted)
			break
		}
	}
	assert.True(t, insertHas, "expected AL_RUWACH (24) to be inserted")

	// Verify update
	updateHas := false
	for _, it := range skillUpdater.ValuesToUpdate {
		if it.SkillID == 1 {
			updateHas = true
			assert.Equal(t, "NV_BASIC", it.Constant.String)
			assert.True(t, it.Constant.Valid)
			assert.False(t, it.Deleted)
			break
		}
	}
	assert.True(t, updateHas, "expected NV_BASIC (1) to be updated")

	// Verify deletion
	if deleted, ok := skillUpdater.DirtyValues[9999]; assert.True(t, ok, "expected 9999 to be edited for deletion") {
		assert.True(t, deleted.Deleted)
	}
}

func TestSkillLoader_loadSkillInfos(t *testing.T) {
	update := domain.Update{
		Date: time.Date(2012, time.January, 2, 0, 0, 0, 0, time.UTC),
		Changes: []domain.UpdateChange{
			{Patch: "testfiles", File: "data/luafiles514/lua files/skillinfoz/JobInheritList.lub"},
			{Patch: "testfiles", File: "data/luafiles514/lua files/skillinfoz/SkillID.lub"},
			{Patch: "testfiles", File: "data/luafiles514/lua files/skillinfoz/SkillInfoList.lub"},
		},
	}

	// Build job and skill tables first (mirrors LoadPatch order)
	jobUpdater := loaders.NewUpdater([]domain.SkillJob{})
	skillUpdater := loaders.NewUpdater([]domain.Skill{
		{SkillID: 24, Constant: domain.NewNullableString("AL_RUWACH"), Description: domain.NewNullableString("Test"), Deleted: false},
	})
	loader := &SkillLoader{}

	loader.loadJobs(testfiles.Root, update, jobUpdater)
	loader.loadSkillIDs(testfiles.Root, update, skillUpdater)

	// Precondition sanity: skill 45 exists from SkillID.lub and is not deleted yet
	if s, ok := skillUpdater.GetForRead(45); assert.True(t, ok) {
		assert.Equal(t, int32(45), s.SkillID)
		assert.False(t, s.Deleted)
	}

	// Execute
	loader.loadSkillInfos(testfiles.Root, update, skillUpdater, jobUpdater)

	// Skills present in SKILL_INFO_LIST should not be marked deleted
	ruwach, ok := skillUpdater.GetForRead(24)
	assert.True(t, ok)
	assert.False(t, ruwach.Deleted)
	assert.Equal(t, "AL_RUWACH", ruwach.Constant.String)
	assert.Equal(t, int32(24), ruwach.SkillID)
	assert.Equal(t, "Ruach", ruwach.Name.String)
	assert.Equal(t, "Test", ruwach.Description.String)
	assert.Equal(t, int32(1), ruwach.MaxLevel.Int32)
	assert.EqualValues(t, []int32{10}, ruwach.SpCost)
	assert.Equal(t, domain.NewNullableBool(false), ruwach.CanSelectLevel)
	assert.EqualValues(t, []int32{10}, ruwach.AttackRange)
	assert.Equal(t, 0, len(ruwach.RequiredSkills))
	assert.Equal(t, 0, len(ruwach.JobRequiredSkills))
}

func TestSkillLoader_loadSkillDesc(t *testing.T) {
	update := domain.Update{
		Date: time.Date(2012, time.January, 2, 0, 0, 0, 0, time.UTC),
		Changes: []domain.UpdateChange{
			{Patch: "testfiles", File: "data/luafiles514/lua files/skillinfoz/SkillDescript.lub"},
		},
	}

	// Seed existing skills so the parser can resolve SKID constants
	existing := []domain.Skill{
		{SkillID: 1, Constant: domain.NewNullableString("NV_BASIC"), Description: domain.NewNullableString("Old"), Deleted: false},
		{SkillID: 24, Constant: domain.NewNullableString("AL_RUWACH"), Description: domain.NewNullableString(""), Deleted: false},
		{SkillID: 45, Constant: domain.NewNullableString("AC_CONCENTRATION"), Description: domain.NewNullableString("something"), Deleted: false},
		{SkillID: 9999, Constant: domain.NewNullableString("TO_DELETE"), Description: domain.NewNullableString("Removed desc"), Deleted: false},
	}
	skillUpdater := loaders.NewUpdater(existing)

	loader := &SkillLoader{}
	loader.loadSkillDesc(testfiles.Root, update, skillUpdater)

	// NV_BASIC description updated
	if s, ok := skillUpdater.GetForRead(1); assert.True(t, ok, "NV_BASIC should exist") {
		assert.False(t, s.Deleted)
		assert.True(t, s.Description.Valid)
		assert.Equal(t, "Basic Skill\nMAX Lv : 9", s.Description.String)
	}

	// AL_RUWACH description updated with 3 lines
	if s, ok := skillUpdater.GetForRead(24); assert.True(t, ok, "AL_RUWACH should exist") {
		assert.False(t, s.Deleted)
		assert.True(t, s.Description.Valid)
		assert.Equal(t, "Ruwach\nMAX Lv : 10\nType : ^000099Offensive^000000", s.Description.String)
	}

	// AC_CONCENTRATION empty description preserved as empty string (npc-like)
	if s, ok := skillUpdater.GetForRead(45); assert.True(t, ok, "AC_CONCENTRATION should exist") {
		assert.False(t, s.Deleted)
		assert.True(t, s.Description.Valid)
		assert.Equal(t, "", s.Description.String)
	}

	// Extra existing skill not present in file should be marked with null description
	if s, ok := skillUpdater.GetForRead(9999); assert.True(t, ok, "9999 should exist") {
		assert.False(t, s.Deleted)
		assert.False(t, s.Description.Valid)
	}
}
