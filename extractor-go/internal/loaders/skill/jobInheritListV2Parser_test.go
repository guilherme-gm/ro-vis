package skill

import (
	"slices"
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestJobInheritV2Parser(t *testing.T) {
	parser := NewJobInehritListV2Parser()

	result := parser.parseFile(testfiles.GetFilePath("rostructs/JobInheritListV2.lua"))

	assert.GreaterOrEqual(t, len(result), 1, "Expected at least one result, got none")

	noviceIdx := slices.IndexFunc(result, func(item SkillJobId) bool { return item.Constant == "JT_NOVICE" })
	assert.GreaterOrEqual(t, noviceIdx, 0, "JT_NOVICE not found")
	assert.Equal(t, "JT_NOVICE", result[noviceIdx].Constant)
	assert.Equal(t, 0, result[noviceIdx].JobId)
	assert.Equal(t, false, result[noviceIdx].InheritedJob.Valid)
	assert.Equal(t, false, result[noviceIdx].InheritedJob2.Valid)

	swordmanIdx := slices.IndexFunc(result, func(item SkillJobId) bool { return item.Constant == "JT_SWORDMAN" })
	assert.GreaterOrEqual(t, swordmanIdx, 0, "JT_SWORDMAN not found")
	assert.Equal(t, "JT_SWORDMAN", result[swordmanIdx].Constant)
	assert.Equal(t, 1, result[swordmanIdx].JobId)
	assert.Equal(t, int32(0), result[swordmanIdx].InheritedJob.Int32)
	assert.Equal(t, true, result[swordmanIdx].InheritedJob.Valid)

	runeKnightHIdx := slices.IndexFunc(result, func(item SkillJobId) bool { return item.Constant == "JT_RUNE_KNIGHT_H" })
	assert.GreaterOrEqual(t, runeKnightHIdx, 0, "JT_RUNE_KNIGHT_H not found")
	assert.Equal(t, "JT_RUNE_KNIGHT_H", result[runeKnightHIdx].Constant)
	assert.Equal(t, 4060, result[runeKnightHIdx].JobId)
	assert.Equal(t, int32(4008), result[runeKnightHIdx].InheritedJob.Int32)
	assert.Equal(t, true, result[runeKnightHIdx].InheritedJob.Valid)
	assert.Equal(t, int32(4054), result[runeKnightHIdx].InheritedJob2.Int32)
	assert.Equal(t, true, result[runeKnightHIdx].InheritedJob2.Valid)
}
