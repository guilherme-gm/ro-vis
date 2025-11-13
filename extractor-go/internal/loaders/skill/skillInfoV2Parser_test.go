package skill

import (
	"slices"
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/jobParsers"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
	"github.com/guilherme-gm/ro-vis/extractor/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestSkillInfoV2Parser(t *testing.T) {
	// Get needed data from other parsers
	jobInehritListV3Parser := jobParsers.NewJobInehritListV3Parser()
	jobIdList := jobInehritListV3Parser.ParseFile(testfiles.GetFilePath("rostructs/JobInheritListV2.lua"))

	jobIdTable := make(map[string]int)
	for _, v := range jobIdList {
		jobIdTable[v.Constant] = v.JobId
	}

	skillIdParser := NewSkillIdV2Parser()
	skillIdTable := skillIdParser.parseFile(testfiles.GetFilePath("rostructs/SkillIDV2.lua"))

	// Actual test
	parser := NewSkillInfoV2Parser().(*SkillInfoV2Parser)

	result := parser.parseFile(testfiles.GetFilePath("rostructs/SkillInfoV2.lua"), jobIdTable, skillIdTable)

	assert.GreaterOrEqual(t, len(result), 1, "Expected at least one result, got none")

	// Checking windwalk (basic skill)
	windwalkIdx := slices.IndexFunc(result, func(s rostructs.SkillInfoV2) bool { return s.Constant == "SN_WINDWALK" })
	assert.GreaterOrEqual(t, windwalkIdx, 0, "SN_WINDWALK not found")

	windwalk := result[windwalkIdx]
	assert.Equal(t, 383, windwalk.SkillId, "SN_WINDWALK has invalid SkillID")
	assert.Equal(t, "Wind Walk", windwalk.SkillName)
	assert.Equal(t, 10, windwalk.MaxLv)
	assert.EqualValues(t, []int{46, 52, 58, 64, 70, 76, 82, 88, 94, 100}, windwalk.SpCost)
	assert.Equal(t, true, windwalk.CanSelectLevel)
	assert.EqualValues(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, windwalk.AttackRange)
	assert.Equal(t, 1, len(windwalk.RequiredSkills))
	assert.Equal(t, 0, len(windwalk.JobRequiredSkills))

	// Checking meltdown (contains RequiredSkills -- not job linked)
	meltdownIdx := slices.IndexFunc(result, func(s rostructs.SkillInfoV2) bool { return s.Constant == "WS_MELTDOWN" })
	assert.GreaterOrEqual(t, meltdownIdx, 0, "WS_MELTDOWN not found")

	meltdown := result[meltdownIdx]
	assert.Equal(t, 384, meltdown.SkillId, "WS_MELTDOWN has invalid SkillID")

	assert.Equal(t, 109, meltdown.RequiredSkills[0].SkillId, "WS_MELTDOWN has invalid first pre-req SkillID")
	assert.Equal(t, 3, meltdown.RequiredSkills[0].Lv, "WS_MELTDOWN has invalid first pre-req Level")

	assert.Equal(t, 105, meltdown.RequiredSkills[1].SkillId, "WS_MELTDOWN has invalid second pre-req SkillID")
	assert.Equal(t, 1, meltdown.RequiredSkills[1].Lv, "WS_MELTDOWN has invalid second pre-req Level")

	// checking arrow vulcan
	arrowVulcanIdx := slices.IndexFunc(result, func(s rostructs.SkillInfoV2) bool { return s.Constant == "CG_ARROWVULCAN" })
	assert.GreaterOrEqual(t, arrowVulcanIdx, 0, "CG_ARROWVULCAN not found")

	arrowVulcan := result[arrowVulcanIdx]
	assert.Equal(t, 394, arrowVulcan.SkillId, "CG_ARROWVULCAN has invalid SkillID")

	assert.Equal(t, 4020, arrowVulcan.JobRequiredSkills[0].Job, "CG_ARROWVULCAN has invalid JobRequiredSkills[0].Job")
	assert.Equal(t, 46, arrowVulcan.JobRequiredSkills[0].RequiredSkills[0].SkillId, "CG_ARROWVULCAN has invalid JobRequiredSkills[0].RequiredSkills[0].SkillId")
	assert.Equal(t, 5, arrowVulcan.JobRequiredSkills[0].RequiredSkills[0].Lv, "CG_ARROWVULCAN has invalid JobRequiredSkills[0].RequiredSkills[0].Lv")
}
