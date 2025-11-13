package idParser

import (
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestSkillIdV2Parser(t *testing.T) {
	parser := NewSkillIdV2Parser()

	result := parser.ParseFile(testfiles.GetFilePath("rostructs/SkillIDV2.lua"))

	assert.GreaterOrEqual(t, len(result), 1, "Expected at least one result, got none")

	assert.Equal(t, 1, result["NV_BASIC"])
}
