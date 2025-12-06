package loaders

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/questParsers"
)

type QuestLoader struct {
	parsers    []questParsers.QuestParser
	repository *repository.QuestRepository
}

// GetRelevantFiles returns a list of all files that are relevant to this loader's parsers.
// The list is deduplicated to avoid returning the same file path multiple times.
func (l *QuestLoader) GetRelevantFiles() []*regexp.Regexp {
	fileMap := make(map[string]bool)
	var result []*regexp.Regexp

	for _, parser := range l.parsers {
		for _, file := range parser.GetRelevantFiles() {
			if !fileMap[file.String()] {
				fileMap[file.String()] = true
				result = append(result, file)
			}
		}
	}

	return result
}

func NewQuestLoader(server *server.Server) *QuestLoader {
	return &QuestLoader{
		parsers: []questParsers.QuestParser{
			questParsers.NewQuestV1Parser(server),
			// questParsers.QuestV2Parser{}, -- Not implemented (uses v1 instead)
			questParsers.NewQuestV3Parser(server),
			questParsers.NewQuestV4Parser(server),
		},
		repository: server.Repositories.QuestRepository,
	}
}

func (l *QuestLoader) LoadPatch(tx *sql.Tx, basePath string, update domain.Update) {
	fmt.Println("> Decoding...")
	var targetParser questParsers.QuestParser = nil
	for _, parser := range l.parsers {
		if parser.IsUpdateInRange(&update) {
			targetParser = parser
			break
		}
	}

	if targetParser == nil {
		panic("Could not find a parser for Quest patch " + update.Name())
	}

	if !targetParser.HasFiles(&update) {
		fmt.Println("Skipped - No meaningful file")
		return
	}

	fileQuests := targetParser.Parse(basePath, &update)

	fmt.Println("> Fetching current list...")
	currentQuests, err := l.repository.GetCurrent(tx)
	if err != nil {
		panic(err)
	}

	questUpdater := NewUpdater(currentQuests, targetParser.FileVersion())
	idExists := make(map[int32]bool)

	for _, fileQuest := range fileQuests {
		idExists[fileQuest.QuestID] = true
		existingQuest, exists := questUpdater.GetForRead(fileQuest.QuestID)
		if !exists || !existingQuest.Equals(fileQuest) {
			updatedQuest := questUpdater.GetForEdit(fileQuest.QuestID)

			// Patch new object
			updatedQuest.Title = fileQuest.Title
			updatedQuest.Description = fileQuest.Description
			updatedQuest.Summary = fileQuest.Summary
			updatedQuest.OldImage = fileQuest.OldImage
			updatedQuest.IconName = fileQuest.IconName
			updatedQuest.NpcSpr = fileQuest.NpcSpr
			updatedQuest.NpcNavi = fileQuest.NpcNavi
			updatedQuest.NpcPosX = fileQuest.NpcPosX
			updatedQuest.NpcPosY = fileQuest.NpcPosY
			updatedQuest.RewardEXP = fileQuest.RewardEXP
			updatedQuest.RewardJEXP = fileQuest.RewardJEXP
			updatedQuest.RewardItemList = fileQuest.RewardItemList
			updatedQuest.CoolTimeQuest = fileQuest.CoolTimeQuest
		}
	}

	for _, quest := range questUpdater.CurrentValues {
		if _, exists := idExists[quest.QuestID]; !exists && !quest.Deleted {
			// Mark for deletion
			deletedQuest := questUpdater.GetForEdit(quest.QuestID)
			temp := *deletedQuest
			*deletedQuest = domain.NewQuest(temp.QuestID, temp.FileVersion)
			deletedQuest.PreviousHistoryID = temp.PreviousHistoryID
			deletedQuest.Deleted = true
		}
	}

	if len(questUpdater.ValuesToInsert) > 0 {
		fmt.Printf("> Saving new records... (%d records to save)\n", len(questUpdater.ValuesToInsert))
		res, err := l.repository.AddToHistory(tx, update.Name(), questUpdater.ValuesToInsert)
		if err != nil {
			panic(err)
		}

		fmt.Println("\tResult: ", res)
	}

	if len(questUpdater.ValuesToUpdate) > 0 {
		fmt.Printf("> Updating records... (%d records to update)\n", len(questUpdater.ValuesToUpdate))
		res, err := l.repository.AddToHistory(tx, update.Name(), questUpdater.ValuesToUpdate)
		if err != nil {
			panic(err)
		}

		fmt.Println("\tResult: ", res)
	}
}

func (l *QuestLoader) Name() string {
	return "quests"
}
