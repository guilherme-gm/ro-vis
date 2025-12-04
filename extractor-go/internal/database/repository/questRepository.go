package repository

import (
	"context"
	"database/sql"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type QuestRepository struct {
	HistoryBaseRepository[
		domain.Quest,
		dao.BulkInsertQuestHistoryParams,
		*dao.BulkInsertQuestHistoryParams,
		*dao.BulkUpsertQuestParams,
	]
}

// NewQuestRepository creates a new QuestRepository instance
func NewQuestRepository(db database.IDatabase) *QuestRepository {
	repo := &QuestRepository{
		HistoryBaseRepository: HistoryBaseRepository[
			domain.Quest,
			dao.BulkInsertQuestHistoryParams,
			*dao.BulkInsertQuestHistoryParams,
			*dao.BulkUpsertQuestParams,
		]{
			DB:       db,
			BulkSize: 500,
		},
	}
	repo.getCurrentData = repo.getCurrentDataFn
	repo.bulkInsertHistory = repo.bulkInsertQuestHistoryFn
	repo.bulkUpsertRecords = repo.bulkUpsertQuestFn
	repo.getIdsInUpdate = repo.getIdsInUpdateFn
	repo.newBulkParamEntry = func() *dao.BulkInsertQuestHistoryParams { return &dao.BulkInsertQuestHistoryParams{} }
	repo.newRecordParam = func() *dao.BulkUpsertQuestParams { return &dao.BulkUpsertQuestParams{} }
	return repo
}

func (r *QuestRepository) getCurrentDataFn(dao dao.IDAO, ctx context.Context) ([]domain.Quest, error) {
	return toDomainSlice(dao.GetCurrentQuests(ctx))
}

func (r *QuestRepository) bulkInsertQuestHistoryFn(dao dao.IDAO, arg []*dao.BulkInsertQuestHistoryParams) (sql.Result, error) {
	return dao.BulkInsertQuestHistory(context.Background(), arg)
}

func (r *QuestRepository) bulkUpsertQuestFn(dao dao.IDAO, arg []*dao.BulkUpsertQuestParams) (sql.Result, error) {
	return dao.BulkUpsertQuests(context.Background(), arg)
}

func (r *QuestRepository) getIdsInUpdateFn(dao dao.IDAO, update string) ([]IdHistory, error) {
	return toIdHistorySlice(dao.GetQuestsIdsInUpdate(context.Background(), update))
}

func (r *QuestRepository) CountChangesInUpdate(tx *sql.Tx, update string) (int, error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.CountChangedQuestsInUpdate(context.Background(), update)
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *QuestRepository) sqlRecordToDomain(dbFrom dao.PreviousQuestHistoryVw, dbTo dao.QuestHistory, lastUpdate sql.NullString) FromToRecord[domain.Quest] {
	var from *domain.Record[domain.Quest] = nil
	var to *domain.Record[domain.Quest] = nil

	if dbFrom.HistoryID.Valid {
		from = &domain.Record[domain.Quest]{
			Update: dbFrom.Update.String,
			Data:   dbFrom.ToDomain(),
		}
	}

	if dbTo.HistoryID != 0 {
		to = &domain.Record[domain.Quest]{
			Update: dbTo.Update,
			Data:   dbTo.ToDomain(),
		}
	}

	return FromToRecord[domain.Quest]{
		LastUpdate: domain.NullableString(lastUpdate),
		From:       from,
		To:         to,
	}
}

func (r *QuestRepository) GetChangesInUpdate(tx *sql.Tx, update string, pagination Pagination) ([]FromToRecord[domain.Quest], error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.GetChangedQuests(context.Background(), dao.GetChangedQuestsParams{
		Update: update,
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []FromToRecord[domain.Quest]{}, nil
	}

	if err != nil {
		return nil, err
	}

	records := make([]FromToRecord[domain.Quest], len(res))
	for idx, qmodel := range res {
		records[idx] = r.sqlRecordToDomain(qmodel.PreviousQuestHistoryVw, qmodel.QuestHistory, qmodel.Lastupdate)
	}

	return records, nil
}

func (r *QuestRepository) GetQuestHistory(tx *sql.Tx, questId int32, pagination Pagination) ([]FromToRecord[domain.Quest], error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.GetQuestHistory(context.Background(), dao.GetQuestHistoryParams{
		QuestID: questId,
		Offset:  pagination.Offset,
		Limit:   pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []FromToRecord[domain.Quest]{}, nil
	}

	if err != nil {
		return nil, err
	}

	records := make([]FromToRecord[domain.Quest], len(res))
	for idx, qmodel := range res {
		records[idx] = r.sqlRecordToDomain(qmodel.PreviousQuestHistoryVw, qmodel.QuestHistory, sql.NullString{})
	}

	return records, nil
}

func (r *QuestRepository) CountQuests(tx *sql.Tx) (int32, error) {
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

func (r *QuestRepository) GetQuests(tx *sql.Tx, pagination Pagination) ([]domain.MinQuest, error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.GetQuestList(context.Background(), dao.GetQuestListParams{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []domain.MinQuest{}, nil
	}
	if err != nil {
		return []domain.MinQuest{}, nil
	}

	quests := make([]domain.MinQuest, len(res))
	for idx, val := range res {
		quests[idx] = domain.MinQuest{
			QuestID:    val.QuestID,
			LastUpdate: val.Lastupdate,
			Title:      domain.NullableString(val.Title),
		}
	}
	return quests, nil
}
