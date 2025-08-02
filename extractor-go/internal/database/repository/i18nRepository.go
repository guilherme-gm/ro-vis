package repository

import (
	"context"
	"database/sql"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type I18nRepository struct {
	DB *database.Database
}

// NewI18nRepository creates a new I18nRepository instance
func NewI18nRepository(db *database.Database) *I18nRepository {
	return &I18nRepository{
		DB: db,
	}
}

func (r *I18nRepository) GetCurrentI18ns(tx *sql.Tx) (*[]domain.I18n, error) {
	queries := r.DB.GetQueries(tx)
	res, err := queries.GetCurrentI18ns(context.Background())
	if err == sql.ErrNoRows {
		return &[]domain.I18n{}, nil
	}

	if err != nil {
		return nil, err
	}

	i18ns := make([]domain.I18n, len(res))
	for idx, qmodel := range res {
		i18ns[idx] = qmodel.ToDomain()
	}

	return &i18ns, nil
}

func (r *I18nRepository) addI18nsToHistory_sub(tx *sql.Tx, update string, newHistories *[]domain.I18n) error {
	queries := r.DB.GetQueries(tx)
	bulkParams := []dao.BulkInsertI18nHistoryParams{}
	updatedIdMap := make(map[uint64]bool, len((*newHistories)))
	for _, it := range *newHistories {
		updatedIdMap[it.I18nId] = true
		bulkParams = append(bulkParams, dao.BulkInsertI18nHistoryParams{
			PreviousHistoryID: sql.NullInt64(it.PreviousHistoryID),
			I18nID:            it.I18nId,
			FileVersion:       it.FileVersion,
			Update:            update,
			ContainerFile:     it.ContainerFile,
			EnText:            it.EnText,
			PtBrText:          it.PtBrText,
			Active:            it.Active,
		})
	}

	_, err := queries.BulkInsertI18nHistory(context.Background(), bulkParams)
	if err != nil {
		return err
	}

	res, err := queries.GetI18nsIdsInUpdate(context.Background(), update)
	if err != nil {
		return err
	}

	upsertParams := []dao.BulkUpsertI18nParams{}
	for _, id := range res {
		if _, ok := updatedIdMap[id.I18nID]; !ok {
			continue
		}

		upsertParams = append(upsertParams, dao.BulkUpsertI18nParams{
			I18nID:    id.I18nID,
			HistoryID: id.HistoryID,
			Deleted:   false,
		})
	}

	_, err = queries.BulkUpsertI18ns(context.Background(), upsertParams)
	if err != nil {
		return err
	}

	return err
}

func (r *I18nRepository) AddI18nsToHistory(tx *sql.Tx, update string, newHistories *[]domain.I18n) error {
	if len(*newHistories) == 0 {
		return nil
	}

	// _, err := r.DB.Connection.Exec("USE `rovis-latam`;")
	// if err != nil {
	// 	return err
	// }

	pageSize := 500
	steps := (len(*newHistories) / pageSize)

	i := 0
	for ; i < steps; i++ {
		slice := (*newHistories)[i*pageSize : (i+1)*pageSize]
		if err := r.addI18nsToHistory_sub(tx, update, &slice); err != nil {
			return err
		}
	}

	slice := (*newHistories)[i*pageSize : len(*newHistories)]
	if err := r.addI18nsToHistory_sub(tx, update, &slice); err != nil {
		return err
	}

	return nil
}

func (r *I18nRepository) AddDeletedI18n(tx *sql.Tx, update string, i18n *domain.I18n) error {
	queries := r.DB.GetQueries(tx)
	res, err := queries.BulkInsertI18nHistory(context.Background(), []dao.BulkInsertI18nHistoryParams{{
		PreviousHistoryID: sql.NullInt64(i18n.HistoryID),
		I18nID:            i18n.I18nId,
		FileVersion:       i18n.FileVersion,
		Update:            update,
	}})

	if err != nil {
		return err
	}

	historyId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	i18n.HistoryID = dao.ToNullableInt64(int64(historyId))

	_, err = queries.UpsertI18n(context.Background(), dao.UpsertI18nParams{
		I18nID:          i18n.I18nId,
		LatestHistoryID: i18n.HistoryID.Int64,
		Deleted:         true,
	})

	return err
}
