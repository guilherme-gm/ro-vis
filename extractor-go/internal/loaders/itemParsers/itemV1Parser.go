package itemParsers

import (
	"database/sql"
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
				newItem := domain.NewItem(int32(entry.ItemID), 1)
				newDB[entry.ItemID] = &newItem
			}

			newDB[entry.ItemID].IdentifiedName = dao.ToNullString(entry.Value)
		}
	}

	// ID#Value# tables
	loadTxtSubTable(basePath, update, newDB, "data/num2itemdisplaynametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry *subparsers.ItemValueTableEntry) {
		if entry == nil {
			item.UnidentifiedName = sql.NullString{}
			return
		}

		item.UnidentifiedName = dao.ToNullString(entry.Value)
	})
	loadTxtSubTable(basePath, update, newDB, "data/idnum2itemresnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry *subparsers.ItemValueTableEntry) {
		if entry == nil {
			item.IdentifiedSprite = sql.NullString{}
			return
		}

		item.IdentifiedSprite = dao.ToNullString(entry.Value)
	})
	loadTxtSubTable(basePath, update, newDB, "data/num2itemresnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry *subparsers.ItemValueTableEntry) {
		if entry == nil {
			item.UnidentifiedSprite = sql.NullString{}
			return
		}

		item.UnidentifiedSprite = dao.ToNullString(entry.Value)
	})
	loadTxtSubTable(basePath, update, newDB, "data/itemslotcounttable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry *subparsers.ItemValueTableEntry) {
		if entry == nil {
			item.SlotCount = 0
			return
		}

		slots, err := strconv.Atoi(entry.Value)
		if err != nil {
			panic("slots for item " + strconv.Itoa(int(entry.ItemID)) + " is not int. (Value: " + entry.Value + ")")
		}

		item.SlotCount = int8(slots)
	})
	loadCardPrefix(basePath, update, newDB)
	loadCardIllustName(basePath, update, newDB)

	// ID#Multiline description# tables
	loadTxtSubTable(basePath, update, newDB, "data/idnum2itemdesctable.txt", subparsers.ParseItemDescTable, func(item *domain.Item, entry *subparsers.ItemDescTableEntry) {
		if entry == nil {
			item.IdentifiedDescription = sql.NullString{}
			return
		}

		item.IdentifiedDescription = dao.ToNullString(entry.Description)
	})
	loadTxtSubTable(basePath, update, newDB, "data/num2itemdesctable.txt", subparsers.ParseItemDescTable, func(item *domain.Item, entry *subparsers.ItemDescTableEntry) {
		if entry == nil {
			item.UnidentifiedDescription = sql.NullString{}
			return
		}

		item.UnidentifiedDescription = dao.ToNullString(entry.Description)
	})

	// ID# tables
	loadBookItems(basePath, update, newDB)
	loadBuyingStoreItems(basePath, update, newDB)
	loadCardPostfix(basePath, update, newDB)

	itemList := make([]domain.Item, len(newDB))
	idx := 0
	for _, v := range newDB {
		itemList[idx] = *v
		idx++
	}

	return itemList
}
