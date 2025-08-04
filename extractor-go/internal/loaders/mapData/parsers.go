package mapData

import (
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

var mapNameTable = "data/mapnametable.txt"
var mapNameTableRegex = regexp.MustCompile(`(?i)^` + mapNameTable + `$`)
var mp3NameTable = "data/mp3nametable.txt"
var mp3NameTableRegex = regexp.MustCompile(`(?i)^` + mp3NameTable + `$`)

var naviLink = "data/luafiles514/lua files/navigation/navi_link_br.lub"
var naviLinkRegex = regexp.MustCompile(`(?i)^` + naviLink + `$`)
var naviMap = "data/luafiles514/lua files/navigation/navi_map_br.lub"
var naviMapRegex = regexp.MustCompile(`(?i)^` + naviMap + `$`)
var naviMob = "data/luafiles514/lua files/navigation/navi_mob_br.lub"
var naviMobRegex = regexp.MustCompile(`(?i)^` + naviMob + `$`)
var naviNpc = "data/luafiles514/lua files/navigation/navi_npc_br.lub"
var naviNpcRegex = regexp.MustCompile(`(?i)^` + naviNpc + `$`)

type NaviMapParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange) []rostructs.NaviMap
}

type NaviNpcParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange) []rostructs.NaviNpc
}

type NaviMobParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange) []rostructs.NaviMob
}

type NaviLinkParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange) []rostructs.NaviLink
}
