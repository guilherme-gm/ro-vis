package repository

import (
	"context"
	"database/sql"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
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

func (r *QuestRepository) GetCurrentQuests() (*[]domain.Quest, error) {
	res, err := r.queries.GetCurrentQuests(context.Background())
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

func (r *QuestRepository) AddQuestToHistory(patch string, fromQuest *domain.Quest, toQuest *domain.Quest) error {
	var previousId sql.NullInt32 = sql.NullInt32{Valid: false}
	if fromQuest != nil {
		previousId = fromQuest.HistoryID
	}

	res, err := r.queries.InsertQuestHistory(context.Background(), dao.InsertQuestHistoryParams{
		PreviousHistoryID: previousId,
		QuestID:           toQuest.QuestID,
		FileVersion:       toQuest.FileVersion,
		Patch:             patch,
		Title:             toQuest.Title,
		Description:       toQuest.Description,
		Summary:           toQuest.Summary,
		OldImage:          toQuest.OldImage,
		IconName:          toQuest.IconName,
		NpcSpr:            toQuest.NpcSpr,
		NpcNavi:           toQuest.NpcNavi,
		NpcPosX:           toQuest.NpcPosX,
		NpcPosY:           toQuest.NpcPosY,
		RewardExp:         toQuest.RewardExp,
		RewardJexp:        toQuest.RewardJexp,
		RewardItemList:    toQuest.RewardItemList,
		CoolTimeQuest:     toQuest.CoolTimeQuest,
	})

	if err != nil {
		return err
	}

	historyId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	toQuest.HistoryID = dao.ToNullInt32(int32(historyId))

	_, err = r.queries.UpsertQuest(context.Background(), dao.UpsertQuestParams{
		QuestID:         toQuest.QuestID,
		LatestHistoryID: toQuest.HistoryID.Int32,
		Deleted:         false,
	})

	return err
}

func (r *QuestRepository) AddDeletedQuest(patch string, quest *domain.Quest) error {
	res, err := r.queries.InsertQuestHistory(context.Background(), dao.InsertQuestHistoryParams{
		PreviousHistoryID: quest.HistoryID,
		QuestID:           quest.QuestID,
		FileVersion:       quest.FileVersion,
		Patch:             patch,
	})

	if err != nil {
		return err
	}

	historyId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	quest.HistoryID = dao.ToNullInt32(int32(historyId))

	_, err = r.queries.UpsertQuest(context.Background(), dao.UpsertQuestParams{
		QuestID:         quest.QuestID,
		LatestHistoryID: quest.HistoryID.Int32,
		Deleted:         true,
	})

	return err
}
