package domain

import "database/sql"

type Quest struct {
	PreviousHistoryID sql.NullInt32
	HistoryID         sql.NullInt32
	QuestID           int32
	FileVersion       int32
	Title             sql.NullString
	Description       sql.NullString
	Summary           sql.NullString
	OldImage          sql.NullString
	IconName          sql.NullString
	NpcSpr            sql.NullString
	NpcNavi           sql.NullString
	NpcPosX           sql.NullInt32
	NpcPosY           sql.NullInt32
	RewardExp         sql.NullString
	RewardJexp        sql.NullString
	RewardItemList    sql.NullString
	CoolTimeQuest     sql.NullInt32
	Deleted           bool
}

func (q *Quest) Equals(otherQuest Quest) bool {
	// FileVersion is not checked, if the file has changed but the quest is the same, we don't care.
	return (q.QuestID == otherQuest.QuestID &&
		q.Title == otherQuest.Title &&
		q.Description == otherQuest.Description &&
		q.Summary == otherQuest.Summary &&
		q.OldImage == otherQuest.OldImage &&
		q.IconName == otherQuest.IconName &&
		q.NpcSpr == otherQuest.NpcSpr &&
		q.NpcNavi == otherQuest.NpcNavi &&
		q.NpcPosX == otherQuest.NpcPosX &&
		q.NpcPosY == otherQuest.NpcPosY &&
		q.RewardExp == otherQuest.RewardExp &&
		q.RewardJexp == otherQuest.RewardJexp &&
		q.RewardItemList == otherQuest.RewardItemList &&
		q.CoolTimeQuest == otherQuest.CoolTimeQuest)
}
