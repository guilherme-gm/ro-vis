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
type NaviNpcV1Parser struct {
	server *server.Server
}

func NewNaviNpcV1Parser(server *server.Server) NaviNpcParser {
	return &NaviNpcV1Parser{
		server: server,
	}
}

func (p NaviNpcV1Parser) IsUpdateInRange(update *domain.Update) bool {
	return true
}

func (p NaviNpcV1Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		naviNpcLatamRegex,
	}
}

func (p NaviNpcV1Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p NaviNpcV1Parser) Parse(basePath string, change *domain.UpdateChange) []rostructs.NaviNpc {
	stringDecoder := decoders.ConvertEucKrToUtf8
	if p.server.Type == server.ServerTypeLATAM {
		stringDecoder = decoders.ConvertWin1252ToUtf8
	}

	var naviNpcs []rostructs.NaviNpc
	result := decoders.DecodeLuaTable(basePath+"/"+change.Patch+"/"+change.File, "Navi_Npc", &naviNpcs, stringDecoder)
	if len(result.NotConsumedPaths) > 0 {
		fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
		panic("Not all keys were consumed.")
	}

	if naviNpcs[len(naviNpcs)-1].MapId == "NULL" {
		return naviNpcs[:len(naviNpcs)-1]
	}

	return naviNpcs
}
