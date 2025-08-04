package mapData

import (
	"fmt"
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

/**
 * First version of NaviMap data
 */
type NaviMobV1Parser struct {
	server *server.Server
}

func NewNaviMobV1Parser(server *server.Server) NaviMobParser {
	return &NaviMobV1Parser{
		server: server,
	}
}

func (p NaviMobV1Parser) IsUpdateInRange(update *domain.Update) bool {
	return true
}

func (p NaviMobV1Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		naviMobRegex,
	}
}

func (p NaviMobV1Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p NaviMobV1Parser) Parse(basePath string, change *domain.UpdateChange) []rostructs.NaviMob {
	stringDecoder := decoders.ConvertEucKrToUtf8
	if p.server.Type == server.ServerTypeLATAM {
		stringDecoder = decoders.ConvertWin1252ToUtf8
	}

	var naviMobs []rostructs.NaviMob
	result := decoders.DecodeLuaTable(basePath+"/"+change.Patch+"/"+change.File, "Navi_Mob", &naviMobs, stringDecoder)
	if len(result.NotConsumedPaths) > 0 {
		fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
		panic("Not all keys were consumed.")
	}

	if naviMobs[len(naviMobs)-1].MapId == "NULL" {
		return naviMobs[:len(naviMobs)-1]
	}

	return naviMobs
}
