package repository

import (
	"context"
	"database/sql"

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
		bulkParams = append(bulkParams, dao.BulkInsertQuestHistoryParams{
			PreviousHistoryID: it.PreviousHistoryID,
			QuestID:           it.QuestID,
			FileVersion:       it.FileVersion,
			Patch:             patch,
			Title:             it.Title,
			Description:       it.Description,
			Summary:           it.Summary,
			OldImage:          it.OldImage,
			IconName:          it.IconName,
			NpcSpr:            it.NpcSpr,
			NpcNavi:           it.NpcNavi,
			NpcPosX:           it.NpcPosX,
			NpcPosY:           it.NpcPosY,
			RewardExp:         it.RewardExp,
			RewardJexp:        it.RewardJexp,
			RewardItemList:    it.RewardItemList,
			CoolTimeQuest:     it.CoolTimeQuest,
		})
	}

	_, err := queries.BulkInsertQuestHistory(context.Background(), bulkParams)
	if err != nil {
		return err
	}

	res, err := queries.GetQuestsIdsInPatch(context.Background(), patch)
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
		PreviousHistoryID: quest.HistoryID,
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

	quest.HistoryID = dao.ToNullInt32(int32(historyId))

	_, err = queries.UpsertQuest(context.Background(), dao.UpsertQuestParams{
		QuestID:         quest.QuestID,
		LatestHistoryID: quest.HistoryID.Int32,
		Deleted:         true,
	})

	return err
}
