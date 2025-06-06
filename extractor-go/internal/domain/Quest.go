package domain

type Quest struct {
	PreviousHistoryID NullableInt32
	HistoryID         NullableInt32
	QuestID           int32
	FileVersion       int32
	Title             NullableString
	Description       NullableString
	Summary           NullableString
	OldImage          NullableString
	IconName          NullableString
	NpcSpr            NullableString
	NpcNavi           NullableString
	NpcPosX           NullableInt32
	NpcPosY           NullableInt32
	RewardEXP         NullableString
	RewardJEXP        NullableString
	RewardItemList    NullableString
	CoolTimeQuest     NullableInt32
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
		q.RewardEXP == otherQuest.RewardEXP &&
		q.RewardJEXP == otherQuest.RewardJEXP &&
		q.RewardItemList == otherQuest.RewardItemList &&
		q.CoolTimeQuest == otherQuest.CoolTimeQuest)
}

type MinQuest struct {
	QuestID    int32
	LastUpdate string
	Title      NullableString
}
