package itemParsers

import (
	"fmt"
	"strings"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

/**
 * 6th Version of Item Data (introduced in 2018-04-18)
 * - System/ItemInfo_true.lub now has the "EffectID" field
 */
type ItemV6Parser struct {
	server *server.Server
}

func NewItemV6Parser(server *server.Server) ItemParser {
	return &ItemV6Parser{
		server: server,
	}
}

func (p ItemV6Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2018, time.April, 17, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(2024, time.March, 11, 0, 0, 0, 0, time.UTC))
}

func (p ItemV6Parser) GetRelevantFiles() []string {
	return []string{
		"data/bookitemnametable.txt",
		"data/buyingstoreitemlist.txt",
		"data/cardpostfixnametable.txt",
		"data/cardprefixnametable.txt",
		"data/num2cardillustnametable.txt",
		"System/itemInfo_true.lub",
		"data/itemmoveinfov5.txt",
	}
}

func (p ItemV6Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p ItemV6Parser) Parse(basePath string, update *domain.Update, existingDB map[int32]*domain.Item) []domain.Item {
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

		itemTbl := []rostructs.ItemV6{}
		result := decoders.DecodeLuaTable(basePath+"/"+change.Patch+"/System/itemInfo_true.lub", "tbl", &itemTbl)
		if len(result.NotConsumedPaths) > 0 {
			fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
			panic("Not all keys were consumed.")
		}

		for _, entry := range itemTbl {
			itemID := int32(entry.ItemID)
			if existingItem, ok := existingDB[itemID]; ok {
				newItem := *existingItem
				newDB[itemID] = &newItem
			} else {
				newItem := domain.NewItem(int32(entry.ItemID), 6)
				newDB[itemID] = &newItem
			}

			newDB[itemID].FileVersion = 6
			newDB[itemID].UnidentifiedName = dao.ToNullableString(entry.UnidentifiedDisplayName)
			newDB[itemID].UnidentifiedSprite = dao.ToNullableString(entry.UnidentifiedResourceName)
			newDB[itemID].UnidentifiedDescription = dao.ToNullableString(strings.Join(entry.UnidentifiedDescriptionName, "\n"))
			newDB[itemID].IdentifiedName = dao.ToNullableString(entry.IdentifiedDisplayName)
			newDB[itemID].IdentifiedSprite = dao.ToNullableString(entry.IdentifiedResourceName)
			newDB[itemID].IdentifiedDescription = dao.ToNullableString(strings.Join(entry.IdentifiedDescriptionName, "\n"))
			newDB[itemID].SlotCount = int8(entry.SlotCount)
			newDB[itemID].ClassNum = dao.ToNullableInt32(int32(entry.ClassNum))
			newDB[itemID].IsCostume = entry.Costume
			newDB[itemID].EffectID = int32(entry.EffectID)
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
