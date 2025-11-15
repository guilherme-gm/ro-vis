package delayParser

import (
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestValidateFileConstats(t *testing.T) {
	err := validateFileConstats(testfiles.GetFilePath("SkillDelayList.lub"), SkFlags)
	assert.NoError(t, err)

	err = validateFileConstats(testfiles.GetFilePath("SkillDelayList.lub"), map[string]int{"SKFLAG_NOREDUCT": 1})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "unknown flags: [SKFLAG_DISABLE_FIXEDCASTING_REDUCTION]")
}
