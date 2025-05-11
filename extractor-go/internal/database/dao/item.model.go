package dao

import (
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func (i *GetCurrentItemsRow) ToDomain() domain.Item {
	var moveInfo domain.ItemMoveInfo
	if i.MoveInfo != nil {
		if err := json.Unmarshal(i.MoveInfo, &moveInfo); err != nil {
			panic(err)
		}
	}

	return domain.Item{
		PreviousHistoryID:       i.PreviousHistoryID,
		HistoryID:               ToNullInt32(i.HistoryID),
		ItemID:                  i.ItemID,
		FileVersion:             i.FileVersion,
		IdentifiedName:          i.IdentifiedName,
		IdentifiedDescription:   i.IdentifiedDescription,
		IdentifiedSprite:        i.IdentifiedSprite,
		UnidentifiedName:        i.UnidentifiedName,
		UnidentifiedDescription: i.UnidentifiedDescription,
		UnidentifiedSprite:      i.UnidentifiedSprite,
		SlotCount:               i.SlotCount,
		IsBook:                  i.IsBook,
		CanUseBuyingStore:       i.CanUseBuyingStore,
		CardPrefix:              i.CardPrefix,
		CardIsPostfix:           i.CardIsPostfix,
		CardIllustration:        i.CardIllustration,
		ClassNum:                i.ClassNum,
		MoveInfo:                moveInfo,
		Deleted:                 i.Deleted,
	}
}
