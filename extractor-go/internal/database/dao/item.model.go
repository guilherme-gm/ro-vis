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
		IsCostume:               i.IsCostume,
		EffectID:                i.EffectID,
		PackageID:               i.PackageID,
		MoveInfo:                moveInfo,
		Deleted:                 i.Deleted,
	}
}

func (i *ItemHistory) ToDomain() domain.Item {
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
		IsCostume:               i.IsCostume,
		EffectID:                i.EffectID,
		PackageID:               i.PackageID,
		MoveInfo:                moveInfo,
	}
}

func (i *PreviousItemHistoryVw) ToDomain() domain.Item {
	var moveInfo domain.ItemMoveInfo
	if i.MoveInfo != nil {
		if err := json.Unmarshal(i.MoveInfo, &moveInfo); err != nil {
			panic(err)
		}
	}

	return domain.Item{
		PreviousHistoryID:       i.PreviousHistoryID,
		HistoryID:               i.HistoryID,
		ItemID:                  i.ItemID.Int32,
		FileVersion:             i.FileVersion.Int32,
		IdentifiedName:          i.IdentifiedName,
		IdentifiedDescription:   i.IdentifiedDescription,
		IdentifiedSprite:        i.IdentifiedSprite,
		UnidentifiedName:        i.UnidentifiedName,
		UnidentifiedDescription: i.UnidentifiedDescription,
		UnidentifiedSprite:      i.UnidentifiedSprite,
		SlotCount:               int8(i.SlotCount.Int16),
		IsBook:                  i.IsBook.Bool,
		CanUseBuyingStore:       i.CanUseBuyingStore.Bool,
		CardPrefix:              i.CardPrefix,
		CardIsPostfix:           i.CardIsPostfix.Bool,
		CardIllustration:        i.CardIllustration,
		ClassNum:                i.ClassNum,
		IsCostume:               i.IsCostume.Bool,
		EffectID:                i.EffectID.Int32,
		PackageID:               i.PackageID.Int32,
		MoveInfo:                moveInfo,
	}
}
