package delayParser

import (
	"slices"
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/idParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
	"github.com/guilherme-gm/ro-vis/extractor/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestSkillDelayV1Parser(t *testing.T) {
	// Get needed data from other parsers
	skillIdParser := idParser.NewSkillIdV2Parser()
	skillIdTable := skillIdParser.ParseFile(testfiles.GetFilePath("rostructs/SkillIDV3.lua"))

	// Actual test
	parser := NewSkillDelayV1Parser().(*SkillDelayV1Parser)

	result := parser.ParseFile(testfiles.GetFilePath("rostructs/SkillDelayListV1.lua"), skillIdTable)

	assert.GreaterOrEqual(t, len(result), 1, "Expected at least one result, got none")

	// Checking arrow vulcan
	arrowVulcanIdx := slices.IndexFunc(result, func(s rostructs.SkillDelayV1) bool { return s.SkillId == 394 })
	assert.GreaterOrEqual(t, arrowVulcanIdx, 0, "CG_ARROWVULCAN not found")

	arrowVulcan := result[arrowVulcanIdx]
	assert.Equal(t, 394, arrowVulcan.SkillId, "CG_ARROWVULCAN has invalid SkillID")
	assert.Equal(t, []int{500, 500, 500, 500, 500, 500, 500, 500, 500, 500}, arrowVulcan.SkillCastFixedDelay)
	assert.Equal(t, []int{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500}, arrowVulcan.SkillCastStatDelay)
	assert.Equal(t, []int{600, 600, 600, 600, 600, 600, 600, 600, 600, 600}, arrowVulcan.SkillGlobalPostDelay)
	assert.Equal(t, []int{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500}, arrowVulcan.SkillSinglePostDelay)

	// Checking charged shout beating
	chargedShoutBeatingIdx := slices.IndexFunc(result, func(s rostructs.SkillDelayV1) bool { return s.SkillId == 10018 })
	assert.GreaterOrEqual(t, chargedShoutBeatingIdx, 0, "GD_CHARGESHOUT_BEATING not found")

	chargedShoutBeating := result[chargedShoutBeatingIdx]
	assert.Equal(t, 10018, chargedShoutBeating.SkillId, "GD_CHARGESHOUT_BEATING has invalid SkillID")
	assert.Equal(t, 2, len(chargedShoutBeating.SkillFlag))
	assert.Equal(t, []int{1, 2}, chargedShoutBeating.SkillFlag)
}
