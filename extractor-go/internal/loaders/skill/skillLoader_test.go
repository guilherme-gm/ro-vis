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
