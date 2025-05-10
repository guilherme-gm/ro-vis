package questParsers

import (
	"fmt"
	"strings"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

/**
 * 2018-03-21 was introduced "System/OnGoingQuestInfoList_True.lub" and "System/RecommendedQuestInfoList_True.lub".
 * This marks the end of "questid2display.txt" and lua files/quest/*.lua
 */
type QuestV3Parser struct{}

func (p QuestV3Parser) IsPatchInRange(patch *domain.Patch) bool {
	return (patch.Date.After(time.Date(2018, time.March, 20, 0, 0, 0, 0, time.UTC)) &&
		patch.Date.Before(time.Date(2020, time.August, 4, 0, 0, 0, 0, time.UTC)))
}

func (p QuestV3Parser) HasFiles(patch *domain.Patch) bool {
	for _, fname := range patch.Files {
		if fname == "System/OngoingQuestInfoList_True.lub" {
			return true
		}

		lowerName := strings.ToLower(fname)
		if lowerName == "system/ongoingquestinfolist_true.lub" {
			fmt.Println("FOUND on lower -- " + fname)
			return true
		}
	}

	return false
}

func (p QuestV3Parser) Parse(patchPath string) []domain.Quest {
	var fileQuests []rostructs.QuestV3
	decoders.DecodeLuaTable(patchPath+"System/OngoingQuestInfoList_True.lub", "QuestInfoList", &fileQuests)

	quests := make([]domain.Quest, len(fileQuests))
	for idx, val := range fileQuests {
		quests[idx] = val.ToDomain()
	}

	return quests
}
