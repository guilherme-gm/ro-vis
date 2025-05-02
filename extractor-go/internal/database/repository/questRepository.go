package repository

import (
	"context"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
)

type QuestRepository struct {
	queries *dao.Queries
}

func newQuestRepository(queries *dao.Queries) *QuestRepository {
	return &QuestRepository{queries: queries}
}

func GetQuestRepository() *QuestRepository {
	if repositoriesCache.questRepository == nil {
		repositoriesCache.questRepository = newQuestRepository(database.GetQueries())
	}

	return repositoriesCache.questRepository
}

func (r *QuestRepository) AddQuestToHistory(quest *dao.QuestHistory) error {
	res, err := r.queries.InsertQuestHistory(context.Background(), dao.InsertQuestHistoryParams{
		PreviousHistoryID: quest.PreviousHistoryID,
		QuestID:           quest.QuestID,
		FileVersion:       quest.FileVersion,
		Patch:             quest.Patch,
		Title:             quest.Title,
		Description:       quest.Description,
		Summary:           quest.Summary,
		OldImage:          quest.OldImage,
		IconName:          quest.IconName,
		NpcSpr:            quest.NpcSpr,
		NpcNavi:           quest.NpcNavi,
		NpcPosX:           quest.NpcPosX,
		NpcPosY:           quest.NpcPosY,
		RewardExp:         quest.RewardExp,
		RewardJexp:        quest.RewardJexp,
		RewardItemList:    quest.RewardItemList,
		CoolTimeQuest:     quest.CoolTimeQuest,
	})

	if err != nil {
		return err
	}

	historyId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	quest.HistoryID = int32(historyId)

	_, err = r.queries.UpsertQuest(context.Background(), dao.UpsertQuestParams{
		QuestID:         quest.QuestID,
		LatestHistoryID: quest.HistoryID,
		Deleted:         false,
	})

	return err
}
