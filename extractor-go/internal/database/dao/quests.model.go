package dao

import (
	"database/sql"
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func (q GetCurrentQuestsRow) ToDomain() domain.Quest {
	var rewardItems []domain.RewardItem
	if q.RewardItemList != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.RewardItemList, &rewardItems)
	}

	return domain.Quest{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         ToNullableInt32(q.HistoryID),
		QuestID:           q.QuestID,
		FileVersion:       q.FileVersion,
		Title:             domain.NullableString(q.Title),
		Description:       domain.NullableString(q.Description),
		Summary:           domain.NullableString(q.Summary),
		OldImage:          domain.NullableString(q.OldImage),
		IconName:          domain.NullableString(q.IconName),
		NpcSpr:            domain.NullableString(q.NpcSpr),
		NpcNavi:           domain.NullableString(q.NpcNavi),
		NpcPosX:           domain.NullableInt32(q.NpcPosX),
		NpcPosY:           domain.NullableInt32(q.NpcPosY),
		RewardEXP:         domain.NullableString(q.RewardExp),
		RewardJEXP:        domain.NullableString(q.RewardJexp),
		RewardItemList:    rewardItems,
		CoolTimeQuest:     domain.NullableInt32(q.CoolTimeQuest),
		Deleted:           q.Deleted,
	}
}

func (q *QuestHistory) ToDomain() domain.Quest {
	var rewardItems []domain.RewardItem
	if q.RewardItemList != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.RewardItemList, &rewardItems)
	}

	return domain.Quest{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         ToNullableInt32(q.HistoryID),
		QuestID:           q.QuestID,
		FileVersion:       q.FileVersion,
		Title:             domain.NullableString(q.Title),
		Description:       domain.NullableString(q.Description),
		Summary:           domain.NullableString(q.Summary),
		OldImage:          domain.NullableString(q.OldImage),
		IconName:          domain.NullableString(q.IconName),
		NpcSpr:            domain.NullableString(q.NpcSpr),
		NpcNavi:           domain.NullableString(q.NpcNavi),
		NpcPosX:           domain.NullableInt32(q.NpcPosX),
		NpcPosY:           domain.NullableInt32(q.NpcPosY),
		RewardEXP:         domain.NullableString(q.RewardExp),
		RewardJEXP:        domain.NullableString(q.RewardJexp),
		RewardItemList:    rewardItems,
		CoolTimeQuest:     domain.NullableInt32(q.CoolTimeQuest),
	}
}

func (q *PreviousQuestHistoryVw) ToDomain() domain.Quest {
	var rewardItems []domain.RewardItem
	if q.RewardItemList != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.RewardItemList, &rewardItems)
	}

	return domain.Quest{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         domain.NullableInt32(q.HistoryID),
		QuestID:           q.QuestID.Int32,
		FileVersion:       q.FileVersion.Int32,
		Title:             domain.NullableString(q.Title),
		Description:       domain.NullableString(q.Description),
		Summary:           domain.NullableString(q.Summary),
		OldImage:          domain.NullableString(q.OldImage),
		IconName:          domain.NullableString(q.IconName),
		NpcSpr:            domain.NullableString(q.NpcSpr),
		NpcNavi:           domain.NullableString(q.NpcNavi),
		NpcPosX:           domain.NullableInt32(q.NpcPosX),
		NpcPosY:           domain.NullableInt32(q.NpcPosY),
		RewardEXP:         domain.NullableString(q.RewardExp),
		RewardJEXP:        domain.NullableString(q.RewardJexp),
		RewardItemList:    rewardItems,
		CoolTimeQuest:     domain.NullableInt32(q.CoolTimeQuest),
	}
}

func (q *BulkInsertQuestHistoryParams) FillFromDomain(quest *domain.Quest, update string) {
	var rewardItemListJson string
	if len(quest.RewardItemList) > 0 {
		jsonBytes, _ := json.Marshal(quest.RewardItemList)
		rewardItemListJson = string(jsonBytes)
	}

	q.PreviousHistoryID = sql.NullInt32(quest.PreviousHistoryID)
	q.QuestID = quest.QuestID
	q.FileVersion = quest.FileVersion
	q.Update = update
	q.Title = sql.NullString(quest.Title)
	q.Description = sql.NullString(quest.Description)
	q.Summary = sql.NullString(quest.Summary)
	q.OldImage = sql.NullString(quest.OldImage)
	q.IconName = sql.NullString(quest.IconName)
	q.NpcSpr = sql.NullString(quest.NpcSpr)
	q.NpcNavi = sql.NullString(quest.NpcNavi)
	q.NpcPosX = sql.NullInt32(quest.NpcPosX)
	q.NpcPosY = sql.NullInt32(quest.NpcPosY)
	q.RewardExp = sql.NullString(quest.RewardEXP)
	q.RewardJexp = sql.NullString(quest.RewardJEXP)
	q.RewardItemList = sql.NullString{String: rewardItemListJson, Valid: len(quest.RewardItemList) > 0}
	q.CoolTimeQuest = sql.NullInt32(quest.CoolTimeQuest)
}

func (q *BulkUpsertQuestParams) Fill(id int32, historyId int32, deleted bool) {
	q.QuestID = id
	q.HistoryID = historyId
	q.Deleted = deleted
}
