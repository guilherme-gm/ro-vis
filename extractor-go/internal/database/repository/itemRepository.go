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
	queries *dao.Queries
}

func newItemRepository(queries *dao.Queries) *ItemRepository {
	return &ItemRepository{queries: queries}
}

func GetItemRepository() *ItemRepository {
	if repositoriesCache.ItemRepository == nil {
		repositoriesCache.ItemRepository = newItemRepository(database.GetQueries())
	}

	return repositoriesCache.ItemRepository
}

func (r *ItemRepository) GetCurrentItems() (*[]domain.Item, error) {
	res, err := r.queries.GetCurrentItems(context.Background())
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

func (r *ItemRepository) addItemsToHistory_sub(update string, newHistories *[]domain.Item) error {
	bulkParams := []dao.BulkInsertItemHistoryParams{}
	updatedIdMap := make(map[int32]bool, len((*newHistories)))
	for _, it := range *newHistories {
		updatedIdMap[it.ItemID] = true
		moveInfoJson, err := json.Marshal(it.MoveInfo)
		if err != nil {
			return err
		}

		bulkParams = append(bulkParams, dao.BulkInsertItemHistoryParams{
			PreviousHistoryID:       it.PreviousHistoryID,
			ItemID:                  it.ItemID,
			FileVersion:             it.FileVersion,
			Update:                  update,
			IdentifiedName:          it.IdentifiedName,
			IdentifiedDescription:   it.IdentifiedDescription,
			IdentifiedSprite:        it.IdentifiedSprite,
			UnidentifiedName:        it.UnidentifiedName,
			UnidentifiedDescription: it.UnidentifiedDescription,
			UnidentifiedSprite:      it.UnidentifiedSprite,
			SlotCount:               it.SlotCount,
			IsBook:                  it.IsBook,
			CanUseBuyingStore:       it.CanUseBuyingStore,
			CardPrefix:              it.CardPrefix,
			CardIsPostfix:           it.CardIsPostfix,
			CardIllustration:        it.CardIllustration,
			ClassNum:                it.ClassNum,
			MoveInfo:                moveInfoJson,
		})
	}

	_, err := r.queries.BulkInsertItemHistory(context.Background(), bulkParams)
	if err != nil {
		return err
	}

	res, err := r.queries.GetItemIdsInUpdate(context.Background(), update)
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

	_, err = r.queries.BulkUpsertItems(context.Background(), upsertParams)
	if err != nil {
		return err
	}

	return err
}

func (r *ItemRepository) AddItemsToHistory(patch string, newHistories *[]domain.Item) error {
	if len(*newHistories) == 0 {
		return nil
	}

	steps := (len(*newHistories) / 500)

	i := 0
	for ; i < steps; i++ {
		slice := (*newHistories)[i*500 : (i+1)*500]
		if err := r.addItemsToHistory_sub(patch, &slice); err != nil {
			return err
		}
	}

	slice := (*newHistories)[i*500 : len(*newHistories)]
	if err := r.addItemsToHistory_sub(patch, &slice); err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) AddDeletedItem(patch string, Item *domain.Item) error {
	res, err := r.queries.BulkInsertItemHistory(context.Background(), []dao.BulkInsertItemHistoryParams{{
		PreviousHistoryID: Item.HistoryID,
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

	Item.HistoryID = dao.ToNullInt32(int32(historyId))

	_, err = r.queries.UpsertItem(context.Background(), dao.UpsertItemParams{
		ItemID:          Item.ItemID,
		LatestHistoryID: Item.HistoryID.Int32,
		Deleted:         true,
	})

	return err
}
