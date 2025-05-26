package itemParsers

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	subparsers "github.com/guilherme-gm/ro-vis/extractor/internal/loaders/itemParsers/subParsers"
)

type txtEntry interface {
	GetItemID() int32
}

func loadTxtSubTable[T txtEntry](basePath string, update *domain.Update, newDB map[int32]*domain.Item, fileName string, parser func(string) (map[int32]T, error), mapper func(*domain.Item, *T)) {
	change, err := update.GetChangeForFile(fileName)
	if err != nil {
		// @TODO: Check if the error was actually the not found one
		return
	}

	itemVal, err := parser(basePath + "/" + change.Patch + "/" + fileName)
	if err != nil {
		panic(err)
	}

	for _, item := range newDB {
		if val, ok := itemVal[item.ItemID]; ok {
			mapper(item, &val)
		} else {
			mapper(item, nil)
		}
	}
}

func loadCardPrefix(basePath string, update *domain.Update, newDB map[int32]*domain.Item) {
	loadTxtSubTable(basePath, update, newDB, "data/cardprefixnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry *subparsers.ItemValueTableEntry) {
		if entry == nil {
			item.CardPrefix = domain.NullableString{}
			return
		}

		item.CardPrefix = dao.ToNullableString(entry.Value)
	})
}

func loadCardIllustName(basePath string, update *domain.Update, newDB map[int32]*domain.Item) {
	loadTxtSubTable(basePath, update, newDB, "data/num2cardillustnametable.txt", subparsers.ParseItemValueTable, func(item *domain.Item, entry *subparsers.ItemValueTableEntry) {
		if entry == nil {
			item.CardIllustration = domain.NullableString{}
			return
		}

		item.CardIllustration = dao.ToNullableString(entry.Value)
	})
}

func loadBookItems(basePath string, update *domain.Update, newDB map[int32]*domain.Item) {
	loadTxtSubTable(basePath, update, newDB, "data/bookitemnametable.txt", subparsers.ParseItemListTable, func(item *domain.Item, entry *subparsers.ItemListEntry) {
		if entry == nil {
			item.IsBook = false
			return
		}

		item.IsBook = true
	})
}

func loadBuyingStoreItems(basePath string, update *domain.Update, newDB map[int32]*domain.Item) {
	loadTxtSubTable(basePath, update, newDB, "data/buyingstoreitemlist.txt", subparsers.ParseItemListTable, func(item *domain.Item, entry *subparsers.ItemListEntry) {
		if entry == nil {
			item.CanUseBuyingStore = false
			return
		}

		item.CanUseBuyingStore = true
	})
}

func loadCardPostfix(basePath string, update *domain.Update, newDB map[int32]*domain.Item) {
	loadTxtSubTable(basePath, update, newDB, "data/cardpostfixnametable.txt", subparsers.ParseItemListTable, func(item *domain.Item, entry *subparsers.ItemListEntry) {
		if entry == nil {
			item.CardIsPostfix = false
			return
		}

		item.CardIsPostfix = true
	})
}

func loadItemMoveInfoV5(basePath string, update *domain.Update, newDB map[int32]*domain.Item) {
	loadTxtSubTable(basePath, update, newDB, "data/itemmoveinfov5.txt", subparsers.ParseItemMoveInfoV5, func(item *domain.Item, entry *subparsers.ItemMoveInfoV5Entry) {
		if entry == nil {
			item.MoveInfo = domain.NewItemMoveInfo()
			return
		}

		moveInfo := domain.NewItemMoveInfo()
		moveInfo.CanDrop = entry.CanDrop
		moveInfo.CanTrade = entry.CanTrade
		moveInfo.CanMoveToStorage = entry.CanMoveToStorage
		moveInfo.CanMoveToCart = entry.CanMoveToCart
		moveInfo.CanSellToNpc = entry.CanSellToNpc
		moveInfo.CanMail = entry.CanMail
		moveInfo.CanAuction = entry.CanAuction
		moveInfo.CanMoveToGuildStorage = entry.CanMoveToGuildStorage
		moveInfo.CommentName = entry.Description

		item.MoveInfo = moveInfo
	})
}
