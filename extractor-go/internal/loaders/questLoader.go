package loaders

import (
	"fmt"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/questParsers"
)

type QuestLoader struct {
	parsers []questParsers.QuestParser
}

func NewQuestLoader() *QuestLoader {
	parserV3 := questParsers.QuestV3Parser{}
	return &QuestLoader{
		parsers: []questParsers.QuestParser{
			&parserV3,
		},
	}
}

func (l *QuestLoader) LoadPatch(patch domain.Patch) {
	patchPath := "../patches/" + patch.Name + "/"

	fmt.Println("> Decoding...")
	var targetParser questParsers.QuestParser = nil
	for _, parser := range l.parsers {
		if parser.IsPatchInRange(&patch) {
			targetParser = parser
			break
		}
	}

	if targetParser == nil {
		panic("Could not find a parser for Quest patch " + patch.Name)
	}

	if !targetParser.HasFiles(&patch) {
		fmt.Println("Skipped - No meaningful file")
		return
	}

	parser := questParsers.QuestV3Parser{}
	fileQuests := parser.Parse(patchPath)

	fmt.Println("> Fetching current list...")
	currentQuests, err := repository.GetQuestRepository().GetCurrentQuests()
	if err != nil {
		panic(err)
	}

	fmt.Println("> Diffing...")
	questMap := make(map[int32]*domain.Quest)
	deletedIds := make(map[int32]bool)
	for _, q := range *currentQuests {
		questMap[q.QuestID] = &q

		if !q.Deleted {
			deletedIds[q.QuestID] = true
		}
	}

	var newQuests []domain.Quest
	var updatedQuests []domain.Quest

	for _, fileQuest := range *fileQuests {
		delete(deletedIds, fileQuest.QuestID)
		existingQuest := questMap[fileQuest.QuestID]
		if existingQuest == nil {
			newQuests = append(newQuests, fileQuest)
			continue
		}

		if !existingQuest.Equals(fileQuest) {
			fileQuest.PreviousHistoryID = existingQuest.HistoryID
			updatedQuests = append(updatedQuests, fileQuest)
		}
	}

	fmt.Printf("> Saving new records... (%d records to save)\n", len(newQuests))
	err = repository.GetQuestRepository().AddQuestsToHistory(patch.Name, &newQuests)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Updating records... (%d records to update)\n", len(updatedQuests))
	err = repository.GetQuestRepository().AddQuestsToHistory(patch.Name, &updatedQuests)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Deleting records... (%d records to delete)\n", len(deletedIds))
	for deletedId := range deletedIds {
		err := repository.GetQuestRepository().AddDeletedQuest(patch.Name, questMap[deletedId])
		if err != nil {
			panic(err)
		}
	}
}
