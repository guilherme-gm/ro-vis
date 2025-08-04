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
type NaviMapV1Parser struct {
	server *server.Server
}

func NewNaviMapV1Parser(server *server.Server) NaviMapParser {
	return &NaviMapV1Parser{
		server: server,
	}
}

func (p NaviMapV1Parser) IsUpdateInRange(update *domain.Update) bool {
	return true
}

func (p NaviMapV1Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		naviMapRegex,
	}
}

func (p NaviMapV1Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p NaviMapV1Parser) Parse(basePath string, change *domain.UpdateChange) []rostructs.NaviMap {
	stringDecoder := decoders.ConvertEucKrToUtf8
	if p.server.Type == server.ServerTypeLATAM {
		stringDecoder = decoders.ConvertWin1252ToUtf8
	}

	var naviMaps []rostructs.NaviMap
	result := decoders.DecodeLuaTable(basePath+"/"+change.Patch+"/"+change.File, "Navi_Map", &naviMaps, stringDecoder)
	if len(result.NotConsumedPaths) > 0 {
		fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
		panic("Not all keys were consumed.")
	}

	if naviMaps[len(naviMaps)-1].MapId == "NULL" {
		return naviMaps[:len(naviMaps)-1]
	}

	return naviMaps
}
