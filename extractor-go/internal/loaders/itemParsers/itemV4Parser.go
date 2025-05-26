package itemParsers

import (
	"strings"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

/**
 * 4th Version of Item Data (introduced in 2017-04-19)
 * - Renamed System/ItemInfo.lub to System/ItemInfo_true.lub
 */
type ItemV4Parser struct{}

func (p ItemV4Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2017, time.April, 18, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(2017, time.August, 1, 0, 0, 0, 0, time.UTC))
}

func (p ItemV4Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles([]string{
		"data/bookitemnametable.txt",
		"data/buyingstoreitemlist.txt",
		"data/cardpostfixnametable.txt",
		"data/cardprefixnametable.txt",
		"data/num2cardillustnametable.txt",
		"System/itemInfo_true.lub",
		"data/itemmoveinfov5.txt",
	})
}

func (p ItemV4Parser) Parse(basePath string, update *domain.Update, existingDB map[int32]*domain.Item) []domain.Item {
	newDB := make(map[int32]*domain.Item, len(existingDB))
	if !update.HasChangedAnyFiles([]string{"System/itemInfo_true.lub"}) {
		for k, v := range existingDB {
			newItem := *v
			newDB[k] = &newItem
		}
	} else {
		change, err := update.GetChangeForFile("System/itemInfo_true.lub")
		if err != nil {
			panic(err)
		}

		itemTbl := []rostructs.ItemV2{}
		decoders.DecodeLuaTable(basePath+"/"+change.Patch+"/System/itemInfo_true.lub", "tbl", &itemTbl)
		for _, entry := range itemTbl {
			itemID := int32(entry.ItemID)
			if existingItem, ok := existingDB[itemID]; ok {
				newItem := *existingItem
				newDB[itemID] = &newItem
			} else {
				newItem := domain.NewItem(int32(entry.ItemID), 4)
				newDB[itemID] = &newItem
			}

			newDB[itemID].FileVersion = 4
			newDB[itemID].UnidentifiedName = dao.ToNullableString(entry.UnidentifiedDisplayName)
			newDB[itemID].UnidentifiedSprite = dao.ToNullableString(entry.UnidentifiedResourceName)
			newDB[itemID].UnidentifiedDescription = dao.ToNullableString(strings.Join(entry.UnidentifiedDescriptionName, "\n"))
			newDB[itemID].IdentifiedName = dao.ToNullableString(entry.IdentifiedDisplayName)
			newDB[itemID].IdentifiedSprite = dao.ToNullableString(entry.IdentifiedResourceName)
			newDB[itemID].IdentifiedDescription = dao.ToNullableString(strings.Join(entry.IdentifiedDescriptionName, "\n"))
			newDB[itemID].SlotCount = int8(entry.SlotCount)
			newDB[itemID].ClassNum = dao.ToNullableInt32(int32(entry.ClassNum))
		}
	}

	loadCardPrefix(basePath, update, newDB)
	loadCardIllustName(basePath, update, newDB)
	loadBookItems(basePath, update, newDB)
	loadBuyingStoreItems(basePath, update, newDB)
	loadCardPostfix(basePath, update, newDB)
	loadItemMoveInfoV5(basePath, update, newDB)

	itemList := make([]domain.Item, len(newDB))
	idx := 0
	for _, v := range newDB {
		itemList[idx] = *v
		idx++
	}

	return itemList
}
