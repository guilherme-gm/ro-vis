package subparsers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
)

type ItemMoveInfoV5Entry struct {
	ItemID                int32
	CanDrop               bool
	CanTrade              bool
	CanMoveToStorage      bool
	CanMoveToCart         bool
	CanSellToNpc          bool
	CanMail               bool
	CanAuction            bool
	CanMoveToGuildStorage bool
	Description           string
}

func (d ItemMoveInfoV5Entry) GetItemID() int32 {
	return d.ItemID
}

func ParseItemMoveInfoV5(filePath string) (map[int32]ItemMoveInfoV5Entry, error) {
	moveInfo := make(map[int32]ItemMoveInfoV5Entry)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return moveInfo, err
	}

	lines := strings.Split(decoders.ConvertToUTF8(string(fileContent)), "\n")
	for idx, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "//") || trimmedLine == "" {
			continue
		}

		var columns []string
		tempColumn := ""
		for _, ch := range trimmedLine {
			if ch == '\t' || ch == ' ' {
				columns = append(columns, tempColumn)
				tempColumn = ""
			} else {
				tempColumn += string(ch)
			}
		}

		columns = append(columns, tempColumn)

		if len(columns) < 9 {
			return moveInfo, fmt.Errorf("moveInfo line %d has %d columns. 9 or 10 expected", (idx + 1), len(columns))
		}

		parsedItemId, err := strconv.Atoi(columns[0])
		if err != nil {
			return moveInfo, err
		}

		itemID := int32(parsedItemId)

		restrictions := []int{}
		for i := 1; i <= 8; i++ {
			if columns[i] == "" {
				// I am not 100% sure on that, the client usually crashes and don't open, but LATAM works with broken files..
				// fmt.Printf("WARN: moveInfo line %d has an empty column %d. Assuming 0\n", (idx + 1), i)
				restrictions = append(restrictions, 0)
				continue
			}

			restriction, err := strconv.Atoi(columns[i])
			if err != nil {
				return moveInfo, err
			}

			if restriction != 0 && restriction != 1 {
				return moveInfo, fmt.Errorf("restriction %d for item %d has an invalid value: %d", i, parsedItemId, restriction)
			}

			restrictions = append(restrictions, restriction)
		}

		description := ""
		if len(columns) == 10 {
			description = strings.TrimPrefix(strings.TrimSpace(columns[9]), "/ ")
		}

		// restriction == 0 means not restricted, thus _can_ do something
		moveInfo[itemID] = ItemMoveInfoV5Entry{
			ItemID:                itemID,
			CanDrop:               restrictions[0] == 0,
			CanTrade:              restrictions[1] == 0,
			CanMoveToStorage:      restrictions[2] == 0,
			CanMoveToCart:         restrictions[3] == 0,
			CanSellToNpc:          restrictions[4] == 0,
			CanMail:               restrictions[5] == 0,
			CanAuction:            restrictions[6] == 0,
			CanMoveToGuildStorage: restrictions[7] == 0,
			Description:           description,
		}
	}

	return moveInfo, nil
}
