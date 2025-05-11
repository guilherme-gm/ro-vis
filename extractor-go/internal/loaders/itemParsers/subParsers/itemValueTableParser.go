package subparsers

import (
	"strconv"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
)

type ItemValueTableEntry struct {
	ItemID int32
	Value  string
}

func (d ItemValueTableEntry) GetItemID() int32 {
	return d.ItemID
}

func ParseItemValueTable(filePath string) ([]ItemValueTableEntry, error) {
	itemValues := []ItemValueTableEntry{}

	lines, err := decoders.DecodeTokenTextTable(filePath, 0)
	if err != nil {
		return itemValues, err
	}

	for i := 0; i < len(lines); i += 2 {
		itemIDStr := lines[i]
		name := lines[i+1]

		itemID, err := strconv.Atoi(itemIDStr)
		if err != nil {
			return itemValues, err
		}

		itemValues = append(itemValues, ItemValueTableEntry{
			ItemID: int32(itemID),
			Value:  name,
		})
	}

	return itemValues, nil
}
