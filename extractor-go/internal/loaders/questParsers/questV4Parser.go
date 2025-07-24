package questParsers

import (
	"fmt"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

/**
 * 2018-03-21 was introduced "System/OnGoingQuestInfoList_True.lub" and "System/RecommendedQuestInfoList_True.lub".
 * This marks the end of "questid2display.txt" and lua files/quest/*.lua
 */
type QuestV4Parser struct{}

func (p QuestV4Parser) IsUpdateInRange(update *domain.Update) bool {
	return (update.Date.After(time.Date(2020, time.August, 4, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)))
}

func (p QuestV4Parser) GetRelevantFiles() []string {
	return []string{
		"System/OngoingQuestInfoList_True.lub",
	}
}

func (p QuestV4Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p QuestV4Parser) Parse(basePath string, update *domain.Update) []domain.Quest {
	change, err := update.GetChangeForFile("System/OngoingQuestInfoList_True.lub")
	if err != nil {
		panic(err)
	}

	var fileQuests []rostructs.QuestV4
	result := decoders.DecodeLuaTable(basePath+"/"+change.Patch+"/System/OngoingQuestInfoList_True.lub", "QuestInfoList", &fileQuests)
	if len(result.NotConsumedPaths) > 0 {
		fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
		panic("Not all keys were consumed.")
	}

	quests := make([]domain.Quest, len(fileQuests))
	for idx, val := range fileQuests {
		quests[idx] = val.ToDomain()
	}

	return quests
}
