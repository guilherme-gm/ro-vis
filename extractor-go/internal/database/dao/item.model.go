package dao

import (
	"database/sql"
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func (i GetCurrentItemsRow) ToDomain() domain.Item {
	var moveInfo domain.ItemMoveInfo
	if i.MoveInfo != nil {
		if err := json.Unmarshal(i.MoveInfo, &moveInfo); err != nil {
			panic(err)
		}
	}

	return domain.Item{
		PreviousHistoryID:       domain.NullableInt32(i.PreviousHistoryID),
		HistoryID:               ToNullableInt32(i.HistoryID),
		ItemID:                  i.ItemID,
		FileVersion:             i.FileVersion,
		IdentifiedName:          domain.NullableString(i.IdentifiedName),
		IdentifiedDescription:   domain.NullableString(i.IdentifiedDescription),
		IdentifiedSprite:        domain.NullableString(i.IdentifiedSprite),
		UnidentifiedName:        domain.NullableString(i.UnidentifiedName),
		UnidentifiedDescription: domain.NullableString(i.UnidentifiedDescription),
		UnidentifiedSprite:      domain.NullableString(i.UnidentifiedSprite),
		SlotCount:               i.SlotCount,
		IsBook:                  i.IsBook,
		CanUseBuyingStore:       i.CanUseBuyingStore,
		CardPrefix:              domain.NullableString(i.CardPrefix),
		CardIsPostfix:           i.CardIsPostfix,
		CardIllustration:        domain.NullableString(i.CardIllustration),
		ClassNum:                domain.NullableInt32(i.ClassNum),
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
		PreviousHistoryID:       domain.NullableInt32(i.PreviousHistoryID),
		HistoryID:               ToNullableInt32(i.HistoryID),
		ItemID:                  i.ItemID,
		FileVersion:             i.FileVersion,
		IdentifiedName:          domain.NullableString(i.IdentifiedName),
		IdentifiedDescription:   domain.NullableString(i.IdentifiedDescription),
		IdentifiedSprite:        domain.NullableString(i.IdentifiedSprite),
		UnidentifiedName:        domain.NullableString(i.UnidentifiedName),
		UnidentifiedDescription: domain.NullableString(i.UnidentifiedDescription),
		UnidentifiedSprite:      domain.NullableString(i.UnidentifiedSprite),
		SlotCount:               i.SlotCount,
		IsBook:                  i.IsBook,
		CanUseBuyingStore:       i.CanUseBuyingStore,
		CardPrefix:              domain.NullableString(i.CardPrefix),
		CardIsPostfix:           i.CardIsPostfix,
		CardIllustration:        domain.NullableString(i.CardIllustration),
		ClassNum:                domain.NullableInt32(i.ClassNum),
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
		PreviousHistoryID:       domain.NullableInt32(i.PreviousHistoryID),
		HistoryID:               domain.NullableInt32(i.HistoryID),
		ItemID:                  i.ItemID.Int32,
		FileVersion:             i.FileVersion.Int32,
		IdentifiedName:          domain.NullableString(i.IdentifiedName),
		IdentifiedDescription:   domain.NullableString(i.IdentifiedDescription),
		IdentifiedSprite:        domain.NullableString(i.IdentifiedSprite),
		UnidentifiedName:        domain.NullableString(i.UnidentifiedName),
		UnidentifiedDescription: domain.NullableString(i.UnidentifiedDescription),
		UnidentifiedSprite:      domain.NullableString(i.UnidentifiedSprite),
		SlotCount:               int8(i.SlotCount.Int16),
		IsBook:                  i.IsBook.Bool,
		CanUseBuyingStore:       i.CanUseBuyingStore.Bool,
		CardPrefix:              domain.NullableString(i.CardPrefix),
		CardIsPostfix:           i.CardIsPostfix.Bool,
		CardIllustration:        domain.NullableString(i.CardIllustration),
		ClassNum:                domain.NullableInt32(i.ClassNum),
		IsCostume:               i.IsCostume.Bool,
		EffectID:                i.EffectID.Int32,
		PackageID:               i.PackageID.Int32,
		MoveInfo:                moveInfo,
	}
}

func (q *BulkInsertItemHistoryParams) FillFromDomain(item *domain.Item, update string) {
	var moveInfoJson []byte
	if item.MoveInfo != (domain.ItemMoveInfo{}) {
		moveInfoJson, _ = json.Marshal(item.MoveInfo)
	}

	q.Update = update
	q.FileVersion = item.FileVersion

	q.PreviousHistoryID = sql.NullInt32(item.PreviousHistoryID)
	q.ItemID = item.ItemID
	q.IdentifiedName = sql.NullString(item.IdentifiedName)
	q.IdentifiedDescription = sql.NullString(item.IdentifiedDescription)
	q.IdentifiedSprite = sql.NullString(item.IdentifiedSprite)
	q.UnidentifiedName = sql.NullString(item.UnidentifiedName)
	q.UnidentifiedDescription = sql.NullString(item.UnidentifiedDescription)
	q.UnidentifiedSprite = sql.NullString(item.UnidentifiedSprite)
	q.SlotCount = int8(item.SlotCount)
	q.IsBook = item.IsBook
	q.CanUseBuyingStore = item.CanUseBuyingStore
	q.CardPrefix = sql.NullString(item.CardPrefix)
	q.CardIsPostfix = item.CardIsPostfix
	q.CardIllustration = sql.NullString(item.CardIllustration)
	q.ClassNum = sql.NullInt32(item.ClassNum)
	q.IsCostume = item.IsCostume
	q.EffectID = item.EffectID
	q.PackageID = item.PackageID
	q.MoveInfo = moveInfoJson
}

func (q *BulkUpsertItemParams) Fill(id int32, historyId int32, deleted bool) {
	q.ItemID = id
	q.HistoryID = historyId
	q.Deleted = deleted
}
