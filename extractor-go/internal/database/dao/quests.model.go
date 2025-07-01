package dao

import (
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func (q *GetCurrentQuestsRow) ToDomain() domain.Quest {
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
