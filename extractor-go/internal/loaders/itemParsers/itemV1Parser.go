package itemParsers

import (
	"strconv"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	subparsers "github.com/guilherme-gm/ro-vis/extractor/internal/loaders/itemParsers/subParsers"
)

/**
 * First version of Item data, everything is TXT
 */
type ItemV1Parser struct{}

func (p ItemV1Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.Before(time.Date(2012, time.July, 11, 0, 0, 0, 0, time.UTC))
}

func (p ItemV1Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles([]string{
		"data/bookitemnametable.txt",
		"data/buyingstoreitemlist.txt",
		"data/cardpostfixnametable.txt",
		"data/cardprefixnametable.txt",
		"data/idnum2itemdesctable.txt",
		"data/idnum2itemdisplaynametable.txt",
		"data/idnum2itemresnametable.txt",
		"data/itemslotcounttable.txt",
		"data/num2cardillustnametable.txt",
		"data/num2itemdesctable.txt",
		"data/num2itemdisplaynametable.txt",
		"data/num2itemresnametable.txt",
	})
}

type txtEntry interface {
	GetItemID() int32
}

func loadTxtSubTable[T txtEntry](basePath string, update *domain.Update, newDB map[int32]*domain.Item, fileName string, parser func(string) ([]T, error), mapper func(*domain.Item, T)) {
	change, err := update.GetChangeForFile(fileName)
	if err != nil {
		// @TODO: Check if the error was actually the not found one
		return
	}

	itemDesc, err := parser(basePath + "/" + change.Patch + "/" + fileName)
	if err != nil {
		panic(err)
	}

	for _, entry := range itemDesc {
		if item, ok := newDB[entry.GetItemID()]; ok {
			mapper(item, entry)
		}
	}
}

func (p ItemV1Parser) Parse(basePath string, update *domain.Update, existingDB map[int32]*domain.Item) []domain.Item {
	newDB := make(map[int32]*domain.Item, len(existingDB))
	if !update.HasChangedAnyFiles([]string{"data/idnum2itemdisplaynametable.txt"}) {
		for k, v := range existingDB {
			newItem := *v
			newDB[k] = &newItem
		}
	} else {
		change, err := update.GetChangeForFile("data/idnum2itemdisplaynametable.txt")
		if err != nil {
			panic(err)
		}

		itemNames, err := subparsers.ParseItemValueTable(basePath + "/" + change.Patch + "/data/idnum2itemdisplaynametable.txt")
		if err != nil {
			panic(err)
		}

		for _, entry := range itemNames {
			if existingItem, ok := existingDB[entry.ItemID]; ok {
				newItem := *existingItem
				newDB[entry.ItemID] = &newItem
			} else {
				newDB[entry.ItemID] = &domain.Item{
					ItemID: int32(entry.ItemID),
				}
			}

			newDB[entry.ItemID].IdentifiedName = dao.ToNullString(entry.Value)
		}
	}

	// ID#Value# tables
	loadTxtSubTable(basePath, update, newDB, "data/num2itemdisplaynametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry subparsers.ItemValueTableEntry) {
		item.UnidentifiedName = dao.ToNullString(entry.Value)
	})
	loadTxtSubTable(basePath, update, newDB, "data/idnum2itemresnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry subparsers.ItemValueTableEntry) {
		item.IdentifiedSprite = dao.ToNullString(entry.Value)
	})
	loadTxtSubTable(basePath, update, newDB, "data/num2itemresnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry subparsers.ItemValueTableEntry) {
		item.UnidentifiedSprite = dao.ToNullString(entry.Value)
	})
	loadTxtSubTable(basePath, update, newDB, "data/cardprefixnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry subparsers.ItemValueTableEntry) {
		item.CardPrefix = dao.ToNullString(entry.Value)
	})
	loadTxtSubTable(basePath, update, newDB, "data/num2cardillustnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry subparsers.ItemValueTableEntry) {
		item.CardIllustration = dao.ToNullString(entry.Value)
	})
	loadTxtSubTable(basePath, update, newDB, "data/itemslotcounttable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry subparsers.ItemValueTableEntry) {
		slots, err := strconv.Atoi(entry.Value)
		if err != nil {
			panic("slots for item " + strconv.Itoa(int(entry.ItemID)) + " is not int. (Value: " + entry.Value + ")")
		}

		item.SlotCount = dao.ToNullInt16(int16(slots))
	})

	// ID#Multiline description# tables
	loadTxtSubTable(basePath, update, newDB, "data/idnum2itemdesctable.txt", subparsers.ParseItemDescTable, func(item *domain.Item, entry subparsers.ItemDescTableEntry) {
		item.IdentifiedDescription = dao.ToNullString(entry.Description)
	})
	loadTxtSubTable(basePath, update, newDB, "data/num2itemdesctable.txt", subparsers.ParseItemDescTable, func(item *domain.Item, entry subparsers.ItemDescTableEntry) {
		item.UnidentifiedDescription = dao.ToNullString(entry.Description)
	})

	// ID# tables
	loadTxtSubTable(basePath, update, newDB, "data/bookitemnametable.txt", subparsers.ParseItemListTable, func(item *domain.Item, entry subparsers.ItemListEntry) {
		item.IsBook = true
	})
	loadTxtSubTable(basePath, update, newDB, "data/buyingstoreitemlist.txt", subparsers.ParseItemListTable, func(item *domain.Item, entry subparsers.ItemListEntry) {
		item.CanUseBuyingStore = true
	})
	loadTxtSubTable(basePath, update, newDB, "data/cardpostfixnametable.txt", subparsers.ParseItemListTable, func(item *domain.Item, entry subparsers.ItemListEntry) {
		item.CardIsPostfix = true
	})

	itemList := make([]domain.Item, len(newDB))
	idx := 0
	for _, v := range newDB {
		itemList[idx] = *v
		idx++
	}

	return itemList
}
