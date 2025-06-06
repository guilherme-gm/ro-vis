package dao

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func (q *GetCurrentQuestsRow) ToDomain() domain.Quest {
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
		RewardItemList:    domain.NullableString(q.RewardItemList),
		CoolTimeQuest:     domain.NullableInt32(q.CoolTimeQuest),
		Deleted:           q.Deleted,
	}
}

func (q *QuestHistory) ToDomain() domain.Quest {
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
		RewardItemList:    domain.NullableString(q.RewardItemList),
		CoolTimeQuest:     domain.NullableInt32(q.CoolTimeQuest),
	}
}

func (q *PreviousQuestHistoryVw) ToDomain() domain.Quest {
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
		RewardItemList:    domain.NullableString(q.RewardItemList),
		CoolTimeQuest:     domain.NullableInt32(q.CoolTimeQuest),
	}
}
