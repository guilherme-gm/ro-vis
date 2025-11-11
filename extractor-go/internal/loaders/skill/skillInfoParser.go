package skill

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

type SkillDescriptParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange, skillIdTable map[string]int) []rostructs.SkillDescript
}

var skillInfoV1Regex = regexp.MustCompile(`(?i)^data/lua files/skillinfo/.*$`)

var skillInfoListV2 = "data/luafiles514/lua files/skillinfoz/skillinfolist.lub"
var skillInfoListV2Regex = regexp.MustCompile(`(?i)^` + skillInfoListV2 + `$`)

var skillIdV2 = "data/luafiles514/lua files/skillinfoz/skillid.lub"
var skillIdV2Regex = regexp.MustCompile(`(?i)^` + skillIdV2 + `$`)

var jobInehritListV2 = "data/luafiles514/lua files/skillinfoz/jobinheritlist.lub"
var jobInehritListV2Regex = regexp.MustCompile(`(?i)^` + jobInehritListV2 + `$`)

var skillDescriptV2 = "data/luafiles514/lua files/skillinfoz/skilldescript.lub"
var skillDescriptV2Regex = regexp.MustCompile(`(?i)^` + skillDescriptV2 + `$`)
