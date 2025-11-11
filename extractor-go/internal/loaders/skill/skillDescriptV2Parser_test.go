package skill

import (
	"slices"
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
	"github.com/guilherme-gm/ro-vis/extractor/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestSkillDescriptV2Parser(t *testing.T) {
	// Get needed data from other parsers
	skillIdParser := NewSkillIdV2Parser()
	skillIdTable := skillIdParser.parseFile(testfiles.GetFilePath("rostructs/SkillIDV2.lua"))

	// Actual test
	parser := NewSkillDescriptV2Parser().(*SkillDescriptV2Parser)

	result := parser.parseFile(testfiles.GetFilePath("rostructs/SkillDescriptV2.lua"), skillIdTable)

	assert.GreaterOrEqual(t, len(result), 1, "Expected at least one result, got none")

	// Checking windwalk (basic skill)
	nvBasicIndex := slices.IndexFunc(result, func(s rostructs.SkillDescript) bool { return s.SkillId == 1 })
	assert.GreaterOrEqual(t, nvBasicIndex, 0, "NV_BASIC (1) not found")

	nvBasic := result[nvBasicIndex]
	assert.Equal(t, 1, nvBasic.SkillId, "NV_BASIC has invalid SkillID")
	assert.EqualValues(t, []string{"Basic Skill", "MAX Lv : 9"}, nvBasic.Description)
}
