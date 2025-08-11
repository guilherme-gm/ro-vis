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
type NaviLinkV1Parser struct {
	server *server.Server
}

func NewNaviLinkV1Parser(server *server.Server) NaviLinkParser {
	return &NaviLinkV1Parser{
		server: server,
	}
}

func (p NaviLinkV1Parser) IsUpdateInRange(update *domain.Update) bool {
	return true
}

func (p NaviLinkV1Parser) GetRelevantFiles() []*regexp.Regexp {
	if p.server.Type == server.ServerTypeLATAM {
		return []*regexp.Regexp{
			regexp.MustCompile(`(?i)^` + naviLinkLatam + `$`),
		}
	} else {
		return []*regexp.Regexp{
			regexp.MustCompile(`(?i)^` + naviLinkKro + `$`),
		}
	}
}

func (p NaviLinkV1Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p NaviLinkV1Parser) Parse(basePath string, change *domain.UpdateChange) []rostructs.NaviLink {
	stringDecoder := decoders.ConvertEucKrToUtf8
	if p.server.Type == server.ServerTypeLATAM {
		stringDecoder = decoders.ConvertWin1252ToUtf8
	}

	var naviLinks []rostructs.NaviLink
	result := decoders.DecodeLuaTable(basePath+"/"+change.Patch+"/"+change.File, "Navi_Link", &naviLinks, stringDecoder)
	if len(result.NotConsumedPaths) > 0 {
		fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
		panic("Not all keys were consumed.")
	}

	if naviLinks[len(naviLinks)-1].MapId == "NULL" {
		return naviLinks[:len(naviLinks)-1]
	}

	return naviLinks
}
