package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type QuestRepository struct{}

func newQuestRepository() *QuestRepository {
	return &QuestRepository{}
}

func GetQuestRepository() *QuestRepository {
	if repositoriesCache.questRepository == nil {
		repositoriesCache.questRepository = newQuestRepository()
	}

	return repositoriesCache.questRepository
}

func (r *QuestRepository) GetCurrentQuests(tx *sql.Tx) (*[]domain.Quest, error) {
	queries := database.GetQueries(tx)
	res, err := queries.GetCurrentQuests(context.Background())
	if err == sql.ErrNoRows {
		return &[]domain.Quest{}, nil
	}

	if err != nil {
		return nil, err
	}

	quests := make([]domain.Quest, len(res))
	for idx, qmodel := range res {
		quests[idx] = qmodel.ToDomain()
	}

	return &quests, nil
}

func (r *QuestRepository) addQuestsToHistory_sub(tx *sql.Tx, patch string, newHistories *[]domain.Quest) error {
	queries := database.GetQueries(tx)
	bulkParams := []dao.BulkInsertQuestHistoryParams{}
	updatedIdMap := make(map[int32]bool, len((*newHistories)))
	for _, it := range *newHistories {
		updatedIdMap[it.QuestID] = true
		var rewardItemListJson string
		if len(it.RewardItemList) > 0 {
			jsonBytes, _ := json.Marshal(it.RewardItemList)
			rewardItemListJson = string(jsonBytes)
		}
		bulkParams = append(bulkParams, dao.BulkInsertQuestHistoryParams{
			PreviousHistoryID: sql.NullInt32(it.PreviousHistoryID),
			QuestID:           it.QuestID,
			FileVersion:       it.FileVersion,
			Patch:             patch,
			Title:             sql.NullString(it.Title),
			Description:       sql.NullString(it.Description),
			Summary:           sql.NullString(it.Summary),
			OldImage:          sql.NullString(it.OldImage),
			IconName:          sql.NullString(it.IconName),
			NpcSpr:            sql.NullString(it.NpcSpr),
			NpcNavi:           sql.NullString(it.NpcNavi),
			NpcPosX:           sql.NullInt32(it.NpcPosX),
			NpcPosY:           sql.NullInt32(it.NpcPosY),
			RewardExp:         sql.NullString(it.RewardEXP),
			RewardJexp:        sql.NullString(it.RewardJEXP),
			RewardItemList:    sql.NullString{String: rewardItemListJson, Valid: len(it.RewardItemList) > 0},
			CoolTimeQuest:     sql.NullInt32(it.CoolTimeQuest),
		})
	}

	_, err := queries.BulkInsertQuestHistory(context.Background(), bulkParams)
	if err != nil {
		return err
	}

	res, err := queries.GetQuestsIdsInUpdate(context.Background(), patch)
	if err != nil {
		return err
	}

	upsertParams := []dao.BulkUpsertQuestParams{}
	for _, id := range res {
		if _, ok := updatedIdMap[id.QuestID]; !ok {
			continue
		}

		upsertParams = append(upsertParams, dao.BulkUpsertQuestParams{
			QuestID:   id.QuestID,
			HistoryID: id.HistoryID,
			Deleted:   false,
		})
	}

	_, err = queries.BulkUpsertQuests(context.Background(), upsertParams)
	if err != nil {
		return err
	}

	return err
}

func (r *QuestRepository) AddQuestsToHistory(tx *sql.Tx, patch string, newHistories *[]domain.Quest) error {
	if len(*newHistories) == 0 {
		return nil
	}

	steps := (len(*newHistories) / 500)

	i := 0
	for ; i < steps; i++ {
		slice := (*newHistories)[i*500 : (i+1)*500]
		if err := r.addQuestsToHistory_sub(tx, patch, &slice); err != nil {
			return err
		}
	}

	slice := (*newHistories)[i*500 : len(*newHistories)]
	if err := r.addQuestsToHistory_sub(tx, patch, &slice); err != nil {
		return err
	}

	return nil
}

func (r *QuestRepository) AddDeletedQuest(tx *sql.Tx, patch string, quest *domain.Quest) error {
	queries := database.GetQueries(tx)
	res, err := queries.BulkInsertQuestHistory(context.Background(), []dao.BulkInsertQuestHistoryParams{{
		PreviousHistoryID: sql.NullInt32(quest.HistoryID),
		QuestID:           quest.QuestID,
		FileVersion:       quest.FileVersion,
		Patch:             patch,
	}})

	if err != nil {
		return err
	}

	historyId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	quest.HistoryID = dao.ToNullableInt32(int32(historyId))

	_, err = queries.UpsertQuest(context.Background(), dao.UpsertQuestParams{
		QuestID:         quest.QuestID,
		LatestHistoryID: quest.HistoryID.Int32,
		Deleted:         true,
	})

	return err
}

func (r *QuestRepository) CountChangesInUpdate(tx *sql.Tx, update string) (int, error) {
	queries := database.GetQueries(tx)
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
	queries := database.GetQueries(tx)
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
	queries := database.GetQueries(tx)
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

func (r *QuestRepository) GetQuests(tx *sql.Tx, pagination Pagination) ([]domain.MinQuest, error) {
	queries := database.GetQueries(tx)
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
