package infoParser

import (
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

type SkillInfoParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange, jobIdTable map[string]int, skillIdTable map[string]int) []rostructs.SkillInfoV2
}

var SkillInfoV1Regex = regexp.MustCompile(`(?i)^data/lua files/skillinfo/.*$`)

var SkillInfoListV2 = "data/luafiles514/lua files/skillinfoz/skillinfolist.lub"
var SkillInfoListV2Regex = regexp.MustCompile(`(?i)^` + SkillInfoListV2 + `$`)
