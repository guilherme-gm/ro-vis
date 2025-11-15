package delayParser

import (
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

type SkillDelayParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange, skillIdTable map[string]int) []rostructs.SkillDelayV1
}

var SkillDelayListV1 = "data/luafiles514/lua files/skillinfoz/skilldelaylist.lub"
var SkillDelayListV1Regex = regexp.MustCompile(`(?i)^` + SkillDelayListV1 + `$`)
