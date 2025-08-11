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

var naviLinkLatam = "data/luafiles514/lua files/navigation/navi_link_br.lub"
var naviLinkKro = "data/luafiles514/lua files/navigation/navi_link_krpri.lub"

var naviMapLatam = "data/luafiles514/lua files/navigation/navi_map_br.lub"
var naviMapKro = "data/luafiles514/lua files/navigation/navi_map_krpri.lub"

var naviMobLatam = "data/luafiles514/lua files/navigation/navi_mob_br.lub"
var naviMobKro = "data/luafiles514/lua files/navigation/navi_mob_krpri.lub"

var naviNpcLatam = "data/luafiles514/lua files/navigation/navi_npc_br.lub"
var naviNpcKro = "data/luafiles514/lua files/navigation/navi_npc_krpri.lub"

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
