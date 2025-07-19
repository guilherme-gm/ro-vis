package loaders

import (
	"database/sql"
	"fmt"

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
func (l *QuestLoader) GetRelevantFiles() []string {
	fileMap := make(map[string]bool)
	var result []string

	for _, parser := range l.parsers {
		for _, file := range parser.GetRelevantFiles() {
			if !fileMap[file] {
				fileMap[file] = true
				result = append(result, file)
			}
		}
	}

	return result
}

func NewQuestLoader(server *server.Server) *QuestLoader {
	return &QuestLoader{
		parsers: []questParsers.QuestParser{
			questParsers.QuestV1Parser{},
			// questParsers.QuestV2Parser{}, -- Not implemented (uses v1 instead)
			questParsers.QuestV3Parser{},
			questParsers.QuestV4Parser{},
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
	currentQuests, err := l.repository.GetCurrentQuests(tx)
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

	for _, fileQuest := range fileQuests {
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
	err = l.repository.AddQuestsToHistory(tx, update.Name(), &newQuests)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Updating records... (%d records to update)\n", len(updatedQuests))
	err = l.repository.AddQuestsToHistory(tx, update.Name(), &updatedQuests)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Deleting records... (%d records to delete)\n", len(deletedIds))
	for deletedId := range deletedIds {
		err := l.repository.AddDeletedQuest(tx, update.Name(), questMap[deletedId])
		if err != nil {
			panic(err)
		}
	}
}
