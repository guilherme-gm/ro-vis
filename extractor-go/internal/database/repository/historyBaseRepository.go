package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type UpdaterEntry[K comparable] interface {
	GetId() K
	GetHistoryId() domain.NullableInt32
}

func toDomainSlice[T any, S DomainConvertible[T]](rows []S, err error) ([]T, error) {
	if err != nil {
		return nil, err
	}

	result := make([]T, len(rows))
	for i, row := range rows {
		result[i] = row.ToDomain()
	}
	return result, nil
}

func toIdHistorySlice[T ~struct {
	HistoryID int32
	ID        int32
}](rows []T, err error) ([]IdHistory, error) {
	if err != nil {
		return nil, err
	}

	result := make([]IdHistory, len(rows))
	for i, row := range rows {
		result[i] = IdHistory(row)
	}
	return result, nil
}

type DomainConvertible[T any] interface {
	ToDomain() T
}

type BulkTypeHist[T any] interface {
}

type BulkTypeHistPtr[T any] interface {
	FillFromDomain(domain *T, updates string)
}

type BulkTypeRecordPtr[T any] interface {
	Fill(id int32, historyId int32, deleted bool)
}

type IdHistory struct {
	HistoryID int32
	ID        int32
}

type Loadable[T any] interface {
	UpdaterEntry[int32]
	IsDeleted() bool
}

type HistoryBaseRepository[
	T Loadable[T],
	BkType BulkTypeHist[T],
	BkTypePtr BulkTypeHistPtr[T],
	RecordTypePtr BulkTypeRecordPtr[T],
] struct {
	DB                database.IDatabase
	BulkSize          int
	getCurrentData    func(dao dao.IDAO, ctx context.Context) ([]T, error)
	bulkInsertHistory func(dao dao.IDAO, arg []BkTypePtr) (sql.Result, error)
	getIdsInUpdate    func(dao dao.IDAO, update string) ([]IdHistory, error)
	bulkUpsertRecords func(dao dao.IDAO, arg []RecordTypePtr) (sql.Result, error)

	countChangesInUpdate func(dao dao.IDAO, update string) (int64, error)
	// factories to allocate concrete params for generics
	newBulkParamEntry func() BkTypePtr
	newRecordParam    func() RecordTypePtr
}

func (r *HistoryBaseRepository[T, BkType, BkTypePtr, RecordTypePtr]) GetCurrent(tx *sql.Tx) ([]T, error) {
	queries := r.DB.GetDAO(tx)
	res, err := r.getCurrentData(queries, context.Background())
	if err == sql.ErrNoRows {
		return []T{}, nil
	}

	return res, err
}

func (r *HistoryBaseRepository[T, BkType, BkTypePtr, RecordTypePtr]) addToHistory_sub(tx *sql.Tx, update string, newItems []T) (AddToHistoryResult, error) {
	queries := r.DB.GetDAO(tx)

	// ObjectID -> deleted (true | false)
	// also means that the object was part of this batch
	idToDeleted := make(map[int32]bool, len(newItems))
	bulkParams := make([]BkTypePtr, len(newItems))
	for idx, it := range newItems {
		idToDeleted[it.GetId()] = it.IsDeleted()

		p := r.newBulkParamEntry()
		p.FillFromDomain(&it, update)
		bulkParams[idx] = p
	}

	addResult := AddToHistoryResult{}
	_, err := r.bulkInsertHistory(queries, bulkParams)
	if err != nil {
		return addResult, fmt.Errorf("failed to bulk insert history, %w", err)
	}

	res, err := r.getIdsInUpdate(queries, update)
	if err != nil {
		return addResult, fmt.Errorf("failed to fetch IDs in current update (\"%s\"). Error: %w", update, err)
	}

	if len(res) == 0 {
		return addResult, fmt.Errorf("failed to fetch IDs in current update (\"%s\"): 0 results were returned -- this should never happen", update)
	}

	var upsertParams []RecordTypePtr
	for _, id := range res {
		if deleted, ok := idToDeleted[id.ID]; ok {
			if deleted {
				addResult.DeletedCount++
			} else {
				addResult.UpsertCount++
			}

			rp := r.newRecordParam()
			rp.Fill(id.ID, id.HistoryID, deleted)
			upsertParams = append(upsertParams, rp)
		}
	}

	_, err = r.bulkUpsertRecords(queries, upsertParams)
	if err != nil {
		return addResult, fmt.Errorf("failed to upsert records: %w", err)
	}

	return addResult, nil
}

func (r *HistoryBaseRepository[T, BkType, BkTypePtr, RecordTypePtr]) AddToHistory(tx *sql.Tx, update string, newItems []T) (AddToHistoryResult, error) {
	if len(newItems) == 0 {
		return AddToHistoryResult{}, nil
	}

	finalResult := AddToHistoryResult{}
	steps := (len(newItems) / r.BulkSize)

	i := 0
	for ; i < steps; i++ {
		slice := newItems[i*r.BulkSize : (i+1)*r.BulkSize]

		res, err := r.addToHistory_sub(tx, update, slice)
		if err != nil {
			return AddToHistoryResult{}, err
		}
		finalResult.UpsertCount += res.UpsertCount
		finalResult.DeletedCount += res.DeletedCount
	}

	slice := newItems[i*r.BulkSize:]
	res, err := r.addToHistory_sub(tx, update, slice)
	if err != nil {
		return AddToHistoryResult{}, err
	}
	finalResult.UpsertCount += res.UpsertCount
	finalResult.DeletedCount += res.DeletedCount

	return finalResult, nil
}

func (r *HistoryBaseRepository[T, BkType, BkTypePtr, RecordTypePtr]) CountChangesInUpdate(tx *sql.Tx, update string) (int, error) {
	res, err := r.countChangesInUpdate(r.DB.GetDAO(tx), update)
	if err != nil {
		return 0, err
	}

	return int(res), nil
}
