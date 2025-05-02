package domain

import "database/sql"

type Quest struct {
	HistoryID      sql.NullInt32
	QuestID        int32
	FileVersion    int32
	Title          sql.NullString
	Description    sql.NullString
	Summary        sql.NullString
	OldImage       sql.NullString
	IconName       sql.NullString
	NpcSpr         sql.NullString
	NpcNavi        sql.NullString
	NpcPosX        sql.NullInt32
	NpcPosY        sql.NullInt32
	RewardExp      sql.NullString
	RewardJexp     sql.NullString
	RewardItemList sql.NullString
	CoolTimeQuest  sql.NullInt32
}
