package repository

import (
	"context"
	"database/sql"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type ItemRepository struct {
	HistoryBaseRepository[
		domain.Item,
		dao.BulkInsertItemHistoryParams,
		*dao.BulkInsertItemHistoryParams,
		*dao.BulkUpsertItemParams,
	]
}

// NewItemRepository creates a new ItemRepository instance
func NewItemRepository(db *database.Database) *ItemRepository {
	repo := &ItemRepository{
		HistoryBaseRepository: HistoryBaseRepository[
			domain.Item,
			dao.BulkInsertItemHistoryParams,
			*dao.BulkInsertItemHistoryParams,
			*dao.BulkUpsertItemParams,
		]{
			DB:       db,
			BulkSize: 500,
		},
	}

	repo.getCurrentData = repo.getCurrentItemsFn
	repo.bulkInsertHistory = repo.bulkInsertItemHistoryFn
	repo.bulkUpsertRecords = repo.bulkUpsertItemFn
	repo.getIdsInUpdate = repo.getIdsInUpdateFn
	repo.newBulkParamEntry = func() *dao.BulkInsertItemHistoryParams { return &dao.BulkInsertItemHistoryParams{} }
	repo.newRecordParam = func() *dao.BulkUpsertItemParams { return &dao.BulkUpsertItemParams{} }

	return repo
}

func (r *ItemRepository) getCurrentItemsFn(dao dao.IDAO, ctx context.Context) ([]domain.Item, error) {
	return toDomainSlice(dao.GetCurrentItems(ctx))
}

func (r *ItemRepository) bulkInsertItemHistoryFn(dao dao.IDAO, arg []*dao.BulkInsertItemHistoryParams) (sql.Result, error) {
	return dao.BulkInsertItemHistory(context.Background(), arg)
}

func (r *ItemRepository) bulkUpsertItemFn(dao dao.IDAO, arg []*dao.BulkUpsertItemParams) (sql.Result, error) {
	return dao.BulkUpsertItems(context.Background(), arg)
}

func (r *ItemRepository) getIdsInUpdateFn(dao dao.IDAO, update string) ([]IdHistory, error) {
	return toIdHistorySlice(dao.GetItemIdsInUpdate(context.Background(), update))
}

func (r *ItemRepository) CountChangesInUpdate(tx *sql.Tx, update string) (int, error) {
	queries := r.DB.GetDAO(tx)
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
	queries := r.DB.GetDAO(tx)
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
	queries := r.DB.GetDAO(tx)
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
	queries := r.DB.GetDAO(tx)

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
	queries := r.DB.GetDAO(tx)
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
