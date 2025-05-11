package subparsers

import (
	"errors"
	"strconv"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
)

type ItemDescTableEntry struct {
	ItemID      int32
	Description string
}

func (d ItemDescTableEntry) GetItemID() int32 {
	return d.ItemID
}

func ParseItemDescTable(filePath string) (map[int32]ItemDescTableEntry, error) {
	itemDescriptions := make(map[int32]ItemDescTableEntry)

	lines, err := decoders.DecodeTokenTextTable(filePath, 1)
	if err != nil {
		return itemDescriptions, err
	}

	idx := 0
	for idx < len(lines) {
		itemID, err := strconv.Atoi(lines[idx])
		if err != nil {
			return itemDescriptions, errors.New("'" + lines[idx] + "' is not a number")
		}
		idx++

		desc := ""
		for idx < len(lines) {
			// Description should read all lines up to a number
			if _, err = strconv.Atoi(lines[idx]); err == nil {
				break
			}

			desc += lines[idx] + "\n"
			idx++
		}

		itemDescriptions[int32(itemID)] = ItemDescTableEntry{
			ItemID:      int32(itemID),
			Description: desc,
		}
	}

	return itemDescriptions, nil
}
