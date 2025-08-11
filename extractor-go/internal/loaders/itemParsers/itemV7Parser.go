package itemParsers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

/**
 * 7th Version of Item Data (introduced in 2024-03-11)
 * - System/ItemInfo_true.lub now has the "PackageID" field
 */
type ItemV7Parser struct {
	server *server.Server
}

func NewItemV7Parser(server *server.Server) ItemParser {
	return &ItemV7Parser{
		server: server,
	}
}

func (p ItemV7Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2024, time.March, 10, 0, 0, 0, 0, time.UTC))
}

func (p ItemV7Parser) GetItemInfoPath() string {
	// LATAM is using itemInfo_new instead of itemInfo_true
	// I am still not 100% sure whether to consider it v7 or v6, or something else
	// as LATAM does not have PackageID, but they have the probability files
	if p.server.Type == server.ServerTypeLATAM {
		return "System/itemInfo_new.lub"
	}

	return "System/itemInfo_true.lub"
}

func (p ItemV7Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		bookItemNameTable,
		buyingStoreItemList,
		cardPostfixNameTable,
		cardPrefixNameTable,
		num2CardIllustNameTable,
		regexp.MustCompile("(?i)^" + p.GetItemInfoPath() + "$"),
		itemMoveInfoV5Table,
	}
}

func (p ItemV7Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p ItemV7Parser) checkNotConsumedPaths(result decoders.LuaDecoderResult) {
	if len(result.NotConsumedPaths) == 0 {
		return
	}

	if p.server.Type == server.ServerTypeLATAM && len(result.NotConsumedPaths) == 1 && result.NotConsumedPaths[0] == "tbl/Visual" {
		// LATAM has created "Visual" field by mistake in some patches, so we ignore it
		// this field is not used.
		return
	}

	fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
	panic("Not all keys were consumed.")
}

func (p ItemV7Parser) Parse(basePath string, update *domain.Update, existingDB map[int32]*domain.Item) []domain.Item {
	newDB := make(map[int32]*domain.Item, len(existingDB))
	if !update.HasChangedAnyFiles([]*regexp.Regexp{regexp.MustCompile("(?i)^" + p.GetItemInfoPath() + "$")}) {
		for k, v := range existingDB {
			newItem := *v
			newDB[k] = &newItem
		}
	} else {
		change, err := update.GetChangeForFile(p.GetItemInfoPath())
		if err != nil {
			panic(err)
		}

		itemTbl := []rostructs.ItemV7{}
		stringDecoder := decoders.ConvertEucKrToUtf8
		if p.server.Type == server.ServerTypeLATAM {
			stringDecoder = decoders.ConvertNoop // LATAM already encodes in UTF-8
		}

		result := decoders.DecodeLuaTable(basePath+"/"+change.Patch+"/"+change.File, "tbl", &itemTbl, stringDecoder)
		p.checkNotConsumedPaths(result)

		for _, entry := range itemTbl {
			itemID := int32(entry.ItemID)
			if existingItem, ok := existingDB[itemID]; ok {
				newItem := *existingItem
				newDB[itemID] = &newItem
			} else {
				newItem := domain.NewItem(int32(entry.ItemID), 7)
				newDB[itemID] = &newItem
			}

			newDB[itemID].FileVersion = 7
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
			newDB[itemID].PackageID = int32(entry.PackageID)
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
