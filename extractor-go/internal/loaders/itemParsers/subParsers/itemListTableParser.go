package subparsers

import (
	"strconv"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
)

type ItemListEntry struct {
	ItemID int32
}

func (d ItemListEntry) GetItemID() int32 {
	return d.ItemID
}

func ParseItemListTable(filePath string) (map[int32]ItemListEntry, error) {
	itemList := make(map[int32]ItemListEntry)

	lines, err := decoders.DecodeTokenTextTable(filePath, 0)
	if err != nil {
		return itemList, err
	}

	for _, itemId := range lines {
		itemID, err := strconv.Atoi(itemId)
		if err != nil {
			return itemList, err
		}

		itemList[int32(itemID)] = ItemListEntry{
			ItemID: int32(itemID),
		}
	}

	return itemList, nil
}
