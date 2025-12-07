package descriptParser

import (
	"slices"
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/idParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
	"github.com/guilherme-gm/ro-vis/extractor/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestSkillDescriptV4Parser(t *testing.T) {
	// Get needed data from other parsers
	skillIdParser := idParser.NewSkillIdV2Parser()
	skillIdTable := skillIdParser.ParseFile(testfiles.GetFilePath("rostructs/SkillIDV3.lua"))

	// Actual test
	parser := NewSkillDescriptV4Parser().(*SkillDescriptV4Parser)

	result := parser.ParseFile(testfiles.GetFilePath("rostructs/SkillDescriptV2.lua"), skillIdTable)

	assert.GreaterOrEqual(t, len(result), 1, "Expected at least one result, got none")

	// Checking windwalk (basic skill)
	nvBasicIndex := slices.IndexFunc(result, func(s rostructs.SkillDescript) bool { return s.SkillId == 1 })
	assert.GreaterOrEqual(t, nvBasicIndex, 0, "NV_BASIC (1) not found")

	nvBasic := result[nvBasicIndex]
	assert.Equal(t, 1, nvBasic.SkillId, "NV_BASIC has invalid SkillID")
	assert.EqualValues(t, []string{"Basic Skill", "MAX Lv : 9"}, nvBasic.Description)
}
