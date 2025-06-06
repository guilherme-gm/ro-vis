package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type ItemRepository struct {
}

func newItemRepository() *ItemRepository {
	return &ItemRepository{}
}

func GetItemRepository() *ItemRepository {
	if repositoriesCache.ItemRepository == nil {
		repositoriesCache.ItemRepository = newItemRepository()
	}

	return repositoriesCache.ItemRepository
}

func (r *ItemRepository) GetCurrentItems(tx *sql.Tx) (*[]domain.Item, error) {
	queries := database.GetQueries(tx)
	res, err := queries.GetCurrentItems(context.Background())
	if err == sql.ErrNoRows {
		return &[]domain.Item{}, nil
	}

	if err != nil {
		return nil, err
	}

	Items := make([]domain.Item, len(res))
	for idx, qmodel := range res {
		Items[idx] = qmodel.ToDomain()
	}

	return &Items, nil
}

func (r *ItemRepository) addItemsToHistory_sub(tx *sql.Tx, update string, newHistories *[]domain.Item) error {
	queries := database.GetQueries(tx)
	bulkParams := []dao.BulkInsertItemHistoryParams{}
	updatedIdMap := make(map[int32]bool, len((*newHistories)))
	for _, it := range *newHistories {
		updatedIdMap[it.ItemID] = true
		moveInfoJson, err := json.Marshal(it.MoveInfo)
		if err != nil {
			return err
		}

		bulkParams = append(bulkParams, dao.BulkInsertItemHistoryParams{
			PreviousHistoryID:       sql.NullInt32(it.PreviousHistoryID),
			ItemID:                  it.ItemID,
			FileVersion:             it.FileVersion,
			Update:                  update,
			IdentifiedName:          sql.NullString(it.IdentifiedName),
			IdentifiedDescription:   sql.NullString(it.IdentifiedDescription),
			IdentifiedSprite:        sql.NullString(it.IdentifiedSprite),
			UnidentifiedName:        sql.NullString(it.UnidentifiedName),
			UnidentifiedDescription: sql.NullString(it.UnidentifiedDescription),
			UnidentifiedSprite:      sql.NullString(it.UnidentifiedSprite),
			SlotCount:               it.SlotCount,
			IsBook:                  it.IsBook,
			CanUseBuyingStore:       it.CanUseBuyingStore,
			CardPrefix:              sql.NullString(it.CardPrefix),
			CardIsPostfix:           it.CardIsPostfix,
			CardIllustration:        sql.NullString(it.CardIllustration),
			ClassNum:                sql.NullInt32(it.ClassNum),
			IsCostume:               it.IsCostume,
			EffectID:                it.EffectID,
			PackageID:               it.PackageID,
			MoveInfo:                moveInfoJson,
		})
	}

	_, err := queries.BulkInsertItemHistory(context.Background(), bulkParams)
	if err != nil {
		return err
	}

	res, err := queries.GetItemIdsInUpdate(context.Background(), update)
	if err != nil {
		return err
	}

	upsertParams := []dao.BulkUpsertItemParams{}
	for _, id := range res {
		if _, ok := updatedIdMap[id.ItemID]; !ok {
			continue
		}

		upsertParams = append(upsertParams, dao.BulkUpsertItemParams{
			ItemID:    id.ItemID,
			HistoryID: id.HistoryID,
			Deleted:   false,
		})
	}

	_, err = queries.BulkUpsertItems(context.Background(), upsertParams)
	if err != nil {
		return err
	}

	return err
}

func (r *ItemRepository) AddItemsToHistory(tx *sql.Tx, patch string, newHistories *[]domain.Item) error {
	if len(*newHistories) == 0 {
		return nil
	}

	steps := (len(*newHistories) / 500)

	i := 0
	for ; i < steps; i++ {
		slice := (*newHistories)[i*500 : (i+1)*500]
		if err := r.addItemsToHistory_sub(tx, patch, &slice); err != nil {
			return err
		}
	}

	slice := (*newHistories)[i*500 : len(*newHistories)]
	if err := r.addItemsToHistory_sub(tx, patch, &slice); err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) AddDeletedItem(tx *sql.Tx, patch string, Item *domain.Item) error {
	queries := database.GetQueries(tx)
	res, err := queries.BulkInsertItemHistory(context.Background(), []dao.BulkInsertItemHistoryParams{{
		PreviousHistoryID: sql.NullInt32(Item.HistoryID),
		ItemID:            Item.ItemID,
		FileVersion:       Item.FileVersion,
		Update:            patch,
	}})

	if err != nil {
		return err
	}

	historyId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	Item.HistoryID = dao.ToNullableInt32(int32(historyId))

	_, err = queries.UpsertItem(context.Background(), dao.UpsertItemParams{
		ItemID:          Item.ItemID,
		LatestHistoryID: Item.HistoryID.Int32,
		Deleted:         true,
	})

	return err
}

func (r *ItemRepository) CountChangesInUpdate(tx *sql.Tx, update string) (int, error) {
	queries := database.GetQueries(tx)
	res, err := queries.CountChangedItemsInUpdate(context.Background(), update)
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *ItemRepository) sqlRecordToDomain(dbFrom dao.PreviousItemHistoryVw, dbTo dao.ItemHistory, lastUpdate sql.NullString) FromToRecord[domain.Item] {
	var from *domain.Record[domain.Item] = nil
	var to *domain.Record[domain.Item] = nil

	if dbFrom.HistoryID.Valid {
		from = &domain.Record[domain.Item]{
			Update: dbFrom.Update.String,
			Data:   dbFrom.ToDomain(),
		}
	}

	if dbTo.HistoryID != 0 {
		to = &domain.Record[domain.Item]{
			Update: dbTo.Update,
			Data:   dbTo.ToDomain(),
		}
	}

	return FromToRecord[domain.Item]{
		LastUpdate: domain.NullableString(lastUpdate),
		From:       from,
		To:         to,
	}
}

func (r *ItemRepository) GetChangesInUpdate(tx *sql.Tx, update string, pagination Pagination) ([]FromToRecord[domain.Item], error) {
	queries := database.GetQueries(tx)
	res, err := queries.GetChangedItems(context.Background(), dao.GetChangedItemsParams{
		Update: update,
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []FromToRecord[domain.Item]{}, nil
	}

	if err != nil {
		return nil, err
	}

	records := make([]FromToRecord[domain.Item], len(res))
	for idx, qmodel := range res {
		records[idx] = r.sqlRecordToDomain(qmodel.PreviousItemHistoryVw, qmodel.ItemHistory, qmodel.Lastupdate)
	}

	return records, nil
}

func (r *ItemRepository) GetItemHistory(tx *sql.Tx, itemId int32, pagination Pagination) ([]FromToRecord[domain.Item], error) {
	queries := database.GetQueries(tx)
	res, err := queries.GetItemHistory(context.Background(), dao.GetItemHistoryParams{
		ItemID: itemId,
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []FromToRecord[domain.Item]{}, nil
	}

	if err != nil {
		return nil, err
	}

	records := make([]FromToRecord[domain.Item], len(res))
	for idx, qmodel := range res {
		records[idx] = r.sqlRecordToDomain(qmodel.PreviousItemHistoryVw, qmodel.ItemHistory, sql.NullString{})
	}

	return records, nil
}

func (r *ItemRepository) CountItems(tx *sql.Tx) (int32, error) {
	queries := database.GetQueries(tx)

	res, err := queries.CountItems(context.Background())
	if err == sql.ErrNoRows {
		return int32(res), nil
	}

	if err != nil {
		return 0, err
	}

	return int32(res), nil
}

func (r *ItemRepository) GetItems(tx *sql.Tx, pagination Pagination) ([]domain.MinItem, error) {
	queries := database.GetQueries(tx)
	res, err := queries.GetItemList(context.Background(), dao.GetItemListParams{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []domain.MinItem{}, nil
	}
	if err != nil {
		return []domain.MinItem{}, nil
	}

	items := make([]domain.MinItem, len(res))
	for idx, val := range res {
		items[idx] = domain.MinItem{
			ItemID:         val.ItemID,
			LastUpdate:     val.Lastupdate,
			IdentifiedName: domain.NullableString(val.IdentifiedName),
		}
	}
	return items, nil
}
