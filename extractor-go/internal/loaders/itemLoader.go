package loaders

import (
	"database/sql"
	"fmt"
	"regexp"

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
func (l *ItemLoader) GetRelevantFiles() []*regexp.Regexp {
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

func NewItemLoader(server *server.Server) *ItemLoader {
	return &ItemLoader{
		parsers: []itemParsers.ItemParser{
			itemParsers.NewItemV1Parser(server),
			itemParsers.NewItemV2Parser(server),
			itemParsers.NewItemV3Parser(server),
			itemParsers.NewItemV4Parser(server),
			itemParsers.NewItemV5Parser(server),
			itemParsers.NewItemV6Parser(server),
			itemParsers.NewItemV7Parser(server),
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
	currentItems, err := l.repository.GetCurrent(tx)
	if err != nil {
		panic(err)
	}

	updater := NewUpdater(currentItems, targetParser.FileVersion())

	// ItemParser requires an input of itemMap due to partial file updates
	// This could probably be solved by providing updater as parameter
	// but I don't feel like refactoring every now...
	itemMap := make(map[int32]*domain.Item)
	for _, q := range currentItems {
		itemMap[q.ItemID] = &q
	}

	fileItems := targetParser.Parse(basePath, &update, itemMap)

	fmt.Println("> Diffing...")

	idExists := make(map[int32]bool)
	for _, fileItem := range fileItems {
		idExists[fileItem.ItemID] = true
		existingItem, exists := updater.GetForRead(fileItem.ItemID)
		if !exists || !existingItem.Equals(fileItem) || existingItem.Deleted {
			updatedItem := updater.GetForEdit(fileItem.ItemID)
			updatedItem.Deleted = false

			// Patch new object
			updatedItem.IdentifiedName = fileItem.IdentifiedName
			updatedItem.IdentifiedDescription = fileItem.IdentifiedDescription
			updatedItem.IdentifiedSprite = fileItem.IdentifiedSprite
			updatedItem.UnidentifiedName = fileItem.UnidentifiedName
			updatedItem.UnidentifiedDescription = fileItem.UnidentifiedDescription
			updatedItem.UnidentifiedSprite = fileItem.UnidentifiedSprite
			updatedItem.SlotCount = fileItem.SlotCount
			updatedItem.IsBook = fileItem.IsBook
			updatedItem.CanUseBuyingStore = fileItem.CanUseBuyingStore
			updatedItem.CardPrefix = fileItem.CardPrefix
			updatedItem.CardIsPostfix = fileItem.CardIsPostfix
			updatedItem.CardIllustration = fileItem.CardIllustration
			updatedItem.ClassNum = fileItem.ClassNum
			updatedItem.IsCostume = fileItem.IsCostume
			updatedItem.EffectID = fileItem.EffectID
			updatedItem.PackageID = fileItem.PackageID
			updatedItem.MoveInfo = fileItem.MoveInfo
		}
	}

	for _, item := range updater.CurrentValues {
		if _, exists := idExists[item.ItemID]; !exists && !item.Deleted {
			// Mark for deletion
			deletedItem := updater.GetForEdit(item.ItemID)
			temp := *deletedItem

			*deletedItem = domain.NewItem(item.ItemID, temp.FileVersion)
			deletedItem.PreviousHistoryID = temp.PreviousHistoryID
			deletedItem.Deleted = true
		}
	}

	if len(updater.ValuesToInsert) > 0 {
		fmt.Printf("> Saving new records... (%d records to save)\n", len(updater.ValuesToInsert))
		res, err := l.repository.AddToHistory(tx, update.Name(), updater.ValuesToInsert)
		if err != nil {
			panic(err)
		}

		fmt.Println("\tResult: ", res)
	}

	if len(updater.ValuesToUpdate) > 0 {
		fmt.Printf("> Updating records... (%d records to update)\n", len(updater.ValuesToUpdate))
		res, err := l.repository.AddToHistory(tx, update.Name(), updater.ValuesToUpdate)
		if err != nil {
			panic(err)
		}

		fmt.Println("\tResult: ", res)
	}
}

func (l *ItemLoader) Name() string {
	return "items"
}
