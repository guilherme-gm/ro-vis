package loaders

import (
	"database/sql"
	"fmt"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/itemParsers"
)

type ItemLoader struct {
	parsers    []itemParsers.ItemParser
	repository *repository.ItemRepository
}

// GetRelevantFiles returns a list of all files that are relevant to this loader's parsers.
// The list is deduplicated to avoid returning the same file path multiple times.
func (l *ItemLoader) GetRelevantFiles() []string {
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

func NewItemLoader(server *server.Server) *ItemLoader {
	return &ItemLoader{
		parsers: []itemParsers.ItemParser{
			itemParsers.ItemV1Parser{},
			itemParsers.ItemV2Parser{},
			itemParsers.ItemV3Parser{},
			itemParsers.ItemV4Parser{},
			itemParsers.ItemV5Parser{},
			itemParsers.ItemV6Parser{},
			itemParsers.ItemV7Parser{},
		},
		repository: server.Repositories.ItemRepository,
	}
}

func (l *ItemLoader) LoadPatch(tx *sql.Tx, basePath string, update domain.Update) {
	fmt.Println("> Decoding...")
	var targetParser itemParsers.ItemParser = nil
	for _, parser := range l.parsers {
		if parser.IsUpdateInRange(&update) {
			targetParser = parser
			break
		}
	}

	if targetParser == nil {
		panic("Could not find a parser for Item patch " + update.Name())
	}

	if !targetParser.HasFiles(&update) {
		fmt.Println("Skipped - No meaningful file")
		return
	}

	fmt.Println("> Fetching current list...")
	currentItems, err := l.repository.GetCurrentItems(tx)
	if err != nil {
		panic(err)
	}

	itemMap := make(map[int32]*domain.Item)
	idsToBeDeleted := make(map[int32]bool)
	for _, q := range *currentItems {
		itemMap[q.ItemID] = &q

		if !q.Deleted {
			idsToBeDeleted[q.ItemID] = true
		}
	}

	fileItems := targetParser.Parse(basePath, &update, itemMap)

	fmt.Println("> Diffing...")

	var newItems []domain.Item
	var updatedItems []domain.Item

	for _, fileItem := range fileItems {
		delete(idsToBeDeleted, fileItem.ItemID)
		existingItem := itemMap[fileItem.ItemID]
		if existingItem == nil {
			newItems = append(newItems, fileItem)
			continue
		}

		if !existingItem.Equals(fileItem) {
			fileItem.PreviousHistoryID = existingItem.HistoryID
			updatedItems = append(updatedItems, fileItem)
		}
	}

	fmt.Printf("> Saving new records... (%d records to save)\n", len(newItems))
	err = l.repository.AddItemsToHistory(tx, update.Name(), &newItems)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Updating records... (%d records to update)\n", len(updatedItems))
	err = l.repository.AddItemsToHistory(tx, update.Name(), &updatedItems)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Deleting records... (%d records to delete)\n", len(idsToBeDeleted))
	for deletedId := range idsToBeDeleted {
		err := l.repository.AddDeletedItem(tx, update.Name(), itemMap[deletedId])
		if err != nil {
			panic(err)
		}
	}
}
