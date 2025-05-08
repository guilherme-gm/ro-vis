// Extensions for quests.sql.go
package dao

import (
	"context"
	"database/sql"
	"strings"
)

const bulkInsertQuestHistoryStart = `-- name: InsertQuestHistory :execresult
INSERT INTO ` + "`" + `quest_history` + "`" + ` (
	` + "`" + `previous_history_id` + "`" + `, -- 1
	` + "`" + `quest_id` + "`" + `, -- 2
	` + "`" + `file_version` + "`" + `, -- 3
	` + "`" + `patch` + "`" + `, -- 4
	` + "`" + `title` + "`" + `, -- 5
	` + "`" + `description` + "`" + `, -- 6
	` + "`" + `summary` + "`" + `, -- 7
	` + "`" + `old_image` + "`" + `, -- 8
	` + "`" + `icon_name` + "`" + `, -- 9
	` + "`" + `npc_spr` + "`" + `, -- 10
	` + "`" + `npc_navi` + "`" + `, -- 11
	` + "`" + `npc_pos_x` + "`" + `, -- 12
	` + "`" + `npc_pos_y` + "`" + `, -- 13
	` + "`" + `reward_exp` + "`" + `, -- 14
	` + "`" + `reward_jexp` + "`" + `, -- 15
	` + "`" + `reward_item_list` + "`" + `, -- 16
	` + "`" + `cool_time_quest` + "`" + ` -- 17
)
VALUES
`

const bulkInsertQuestHistoryValueBlock = `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),`

type BulkInsertQuestHistoryParams struct {
	PreviousHistoryID sql.NullInt32
	QuestID           int32
	FileVersion       int32
	Patch             string
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
}

func (q *Queries) BulkInsertQuestHistory(ctx context.Context, arg []BulkInsertQuestHistoryParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkInsertQuestHistoryStart
	var params []any
	for _, data := range arg {
		query += bulkInsertQuestHistoryValueBlock
		params = append(params,
			data.PreviousHistoryID,
			data.QuestID,
			data.FileVersion,
			data.Patch,
			data.Title,
			data.Description,
			data.Summary,
			data.OldImage,
			data.IconName,
			data.NpcSpr,
			data.NpcNavi,
			data.NpcPosX,
			data.NpcPosY,
			data.RewardExp,
			data.RewardJexp,
			data.RewardItemList,
			data.CoolTimeQuest)
	}
	query = strings.TrimRight(query, ",")
	return q.db.ExecContext(ctx, query, params...)
}

const bulkUpsertQuestsStart = `-- name: BulkUpsertQuestHistory :execresult
INSERT INTO ` + "`quests` (`quest_id`,`latest_history_id`,`deleted`)" + `
VALUES
`

const bulkUpsertQuestsValue = `(?, ?, ?),`

const bulkUpsertQuestsDuplicate = `
ON DUPLICATE KEY UPDATE
	` + "`quests`.`latest_history_id` = VALUES(`quests`.`latest_history_id`)," + `
	` + "`quests`.`deleted` = VALUES(`quests`.`deleted`)"

type BulkUpsertQuestParams struct {
	QuestID   int32
	HistoryID int32
	Deleted   bool
}

func (q *Queries) BulkUpsertQuests(ctx context.Context, arg []BulkUpsertQuestParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkUpsertQuestsStart
	var params []any
	for _, data := range arg {
		query += bulkUpsertQuestsValue
		params = append(params,
			data.QuestID,
			data.HistoryID,
			data.Deleted)
	}
	query = strings.TrimRight(query, ",")

	query += bulkUpsertQuestsDuplicate

	return q.db.ExecContext(ctx, query, params...)
}
