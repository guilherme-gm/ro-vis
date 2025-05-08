package dao

import "github.com/guilherme-gm/ro-vis/extractor/internal/domain"

func (q *GetCurrentQuestsRow) ToDomain() domain.Quest {
	return domain.Quest{
		PreviousHistoryID: q.PreviousHistoryID,
		HistoryID:         ToNullInt32(q.HistoryID),
		QuestID:           q.QuestID,
		FileVersion:       q.FileVersion,
		Title:             q.Title,
		Description:       q.Description,
		Summary:           q.Summary,
		OldImage:          q.OldImage,
		IconName:          q.IconName,
		NpcSpr:            q.NpcSpr,
		NpcNavi:           q.NpcNavi,
		NpcPosX:           q.NpcPosX,
		NpcPosY:           q.NpcPosY,
		RewardExp:         q.RewardExp,
		RewardJexp:        q.RewardJexp,
		RewardItemList:    q.RewardItemList,
		CoolTimeQuest:     q.CoolTimeQuest,
		Deleted:           q.Deleted,
	}
}
