package itemParsers

import (
	"database/sql"
	"strings"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	subparsers "github.com/guilherme-gm/ro-vis/extractor/internal/loaders/itemParsers/subParsers"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

/**
 * 2nd Version of Item Data (introduced in 2012-07-11)
 * - Introduction of System/itemInfo.lub
 * - Removal of data/idnum2item*, data/num2item* and data/itemslotcounttable.txt TXT files
 */
type ItemV2Parser struct{}

func (p ItemV2Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2012, time.July, 10, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(2015, time.April, 21, 0, 0, 0, 0, time.UTC))
}

func (p ItemV2Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles([]string{
		"data/bookitemnametable.txt",
		"data/buyingstoreitemlist.txt",
		"data/cardpostfixnametable.txt",
		"data/cardprefixnametable.txt",
		"data/num2cardillustnametable.txt",
		"System/itemInfo.lub",
	})
}

func (p ItemV2Parser) Parse(basePath string, update *domain.Update, existingDB map[int32]*domain.Item) []domain.Item {
	newDB := make(map[int32]*domain.Item, len(existingDB))
	if !update.HasChangedAnyFiles([]string{"System/itemInfo.lub"}) {
		for k, v := range existingDB {
			newItem := *v
			newDB[k] = &newItem
		}
	} else {
		change, err := update.GetChangeForFile("System/itemInfo.lub")
		if err != nil {
			panic(err)
		}

		itemTbl := []rostructs.ItemV2{}
		decoders.DecodeLuaTable(basePath+"/"+change.Patch+"/System/itemInfo.lub", "tbl", &itemTbl)
		for _, entry := range itemTbl {
			itemID := int32(entry.ItemID)
			if existingItem, ok := existingDB[itemID]; ok {
				newItem := *existingItem
				newDB[itemID] = &newItem
			} else {
				newItem := domain.NewItem(int32(entry.ItemID), 2)
				newDB[itemID] = &newItem
			}

			newDB[itemID].UnidentifiedName = dao.ToNullString(entry.UnidentifiedDisplayName)
			newDB[itemID].UnidentifiedSprite = dao.ToNullString(entry.UnidentifiedResourceName)
			newDB[itemID].UnidentifiedDescription = dao.ToNullString(strings.Join(entry.UnidentifiedDescriptionName, "\n"))
			newDB[itemID].IdentifiedName = dao.ToNullString(entry.IdentifiedDisplayName)
			newDB[itemID].IdentifiedSprite = dao.ToNullString(entry.IdentifiedResourceName)
			newDB[itemID].IdentifiedDescription = dao.ToNullString(strings.Join(entry.IdentifiedDescriptionName, "\n"))
			newDB[itemID].SlotCount = int8(entry.SlotCount)
			newDB[itemID].ClassNum = dao.ToNullInt32(int32(entry.ClassNum))
		}
	}

	// ID#Value# tables
	loadTxtSubTable(basePath, update, newDB, "data/cardprefixnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry *subparsers.ItemValueTableEntry) {
		if entry == nil {
			item.CardPrefix = sql.NullString{}
			return
		}

		item.CardPrefix = dao.ToNullString(entry.Value)
	})
	loadTxtSubTable(basePath, update, newDB, "data/num2cardillustnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry *subparsers.ItemValueTableEntry) {
		if entry == nil {
			item.CardIllustration = sql.NullString{}
			return
		}

		item.CardIllustration = dao.ToNullString(entry.Value)
	})

	// ID# tables
	loadTxtSubTable(basePath, update, newDB, "data/bookitemnametable.txt", subparsers.ParseItemListTable, func(item *domain.Item, entry *subparsers.ItemListEntry) {
		if entry == nil {
			item.IsBook = false
			return
		}

		item.IsBook = true
	})
	loadTxtSubTable(basePath, update, newDB, "data/buyingstoreitemlist.txt", subparsers.ParseItemListTable, func(item *domain.Item, entry *subparsers.ItemListEntry) {
		if entry == nil {
			item.CanUseBuyingStore = false
			return
		}

		item.CanUseBuyingStore = true
	})
	loadTxtSubTable(basePath, update, newDB, "data/cardpostfixnametable.txt", subparsers.ParseItemListTable, func(item *domain.Item, entry *subparsers.ItemListEntry) {
		if entry == nil {
			item.CardIsPostfix = false
			return
		}

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
