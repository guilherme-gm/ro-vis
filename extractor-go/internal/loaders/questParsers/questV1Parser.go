package questParsers

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

/**
 * This is the first version of the Quest Log system.
 * At some point (I think 2007-08-29), the Quest UI was introduced.
 *
 * It contained only questID2displayname.txt, a TokenTextTable.
 */
type QuestV1Parser struct {
	server *server.Server
}

func NewQuestV1Parser(server *server.Server) QuestV1Parser {
	return QuestV1Parser{server: server}
}

func (p QuestV1Parser) IsUpdateInRange(update *domain.Update) bool {
	/**
	 * Notes:
	 * 1) I don't know the real start date (probably ~2007-08-29), but it is older than RO Vis covers
	 * 2) This is also including QuestV2, because RO Vis is not supporting the extra files (at least for now)
	 */
	return update.Date.Before(time.Date(2018, time.March, 20, 0, 0, 0, 0, time.UTC))
}

func (p QuestV1Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile("(?i)^data/questid2display.txt$"),
	}
}

func (p QuestV1Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p QuestV1Parser) Parse(basePath string, update *domain.Update) []domain.Quest {
	change, err := update.GetChangeForFile("data/questid2display.txt")
	if err != nil {
		panic(err)
	}

	stringList, err := decoders.DecodeTokenTextTable(basePath+"/"+change.Patch+"/data/questid2display.txt", 0)
	if err != nil {
		panic(err)
	}

	var quests []domain.Quest
	for len(stringList) > 6 {
		questId, err := strconv.Atoi(stringList[0])
		if err != nil {
			fmt.Println("Failed to parse QuestID: ", stringList[0], err)
			stringList = stringList[6:]
			continue
		}

		qv1 := rostructs.QuestV1{
			QuestId:     questId,
			Title:       stringList[1],
			Icon:        stringList[2],
			Image:       stringList[3],
			Description: stringList[4],
			Summary:     stringList[5],
		}
		quests = append(quests, qv1.ToDomain())

		stringList = stringList[6:]
	}

	return quests
}
