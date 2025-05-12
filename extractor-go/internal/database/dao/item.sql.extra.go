// Extensions for items.sql.go
package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
)

const bulkInsertItemHistoryStart = `-- name: InsertItemHistory :execresult
INSERT INTO ` + "`" + `item_history` + "`" + ` (
	` + "`" + `previous_history_id` + "`" + `,
	` + "`" + `item_id` + "`" + `,
	` + "`" + `file_version` + "`" + `,
	` + "`" + `update` + "`" + `,
	` + "`" + `identified_name` + "`" + `,
	` + "`" + `identified_description` + "`" + `,
	` + "`" + `identified_sprite` + "`" + `,
	` + "`" + `unidentified_name` + "`" + `,
	` + "`" + `unidentified_description` + "`" + `,
	` + "`" + `unidentified_sprite` + "`" + `,
	` + "`" + `slot_count` + "`" + `,
	` + "`" + `is_book` + "`" + `,
	` + "`" + `can_use_buying_store` + "`" + `,
	` + "`" + `card_prefix` + "`" + `,
	` + "`" + `card_is_postfix` + "`" + `,
	` + "`" + `card_illustration` + "`" + `,
	` + "`" + `class_num` + "`" + `,
	` + "`" + `move_info` + "`" + `
)
VALUES
`

const bulkInsertItemHistoryValueBlock = `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),`

type BulkInsertItemHistoryParams struct {
	PreviousHistoryID       sql.NullInt32
	ItemID                  int32
	FileVersion             int32
	Update                  string
	IdentifiedName          sql.NullString
	IdentifiedDescription   sql.NullString
	IdentifiedSprite        sql.NullString
	UnidentifiedName        sql.NullString
	UnidentifiedDescription sql.NullString
	UnidentifiedSprite      sql.NullString
	SlotCount               int8
	IsBook                  bool
	CanUseBuyingStore       bool
	CardPrefix              sql.NullString
	CardIsPostfix           bool
	CardIllustration        sql.NullString
	ClassNum                sql.NullInt32
	MoveInfo                json.RawMessage
}

func (q *Queries) BulkInsertItemHistory(ctx context.Context, arg []BulkInsertItemHistoryParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkInsertItemHistoryStart
	var params []any
	for _, data := range arg {
		query += bulkInsertItemHistoryValueBlock
		params = append(params,
			data.PreviousHistoryID,
			data.ItemID,
			data.FileVersion,
			data.Update,
			data.IdentifiedName,
			data.IdentifiedDescription,
			data.IdentifiedSprite,
			data.UnidentifiedName,
			data.UnidentifiedDescription,
			data.UnidentifiedSprite,
			data.SlotCount,
			data.IsBook,
			data.CanUseBuyingStore,
			data.CardPrefix,
			data.CardIsPostfix,
			data.CardIllustration,
			data.ClassNum,
			data.MoveInfo)
	}
	query = strings.TrimRight(query, ",")
	return q.db.ExecContext(ctx, query, params...)
}

const bulkUpsertItemsStart = `-- name: BulkUpsertItemHistory :execresult
INSERT INTO ` + "`items` (`item_id`,`latest_history_id`,`deleted`)" + `
VALUES
`

const bulkUpsertItemsValue = `(?, ?, ?),`

const bulkUpsertItemsDuplicate = `
ON DUPLICATE KEY UPDATE
	` + "`items`.`latest_history_id` = VALUES(`items`.`latest_history_id`)," + `
	` + "`items`.`deleted` = VALUES(`items`.`deleted`)"

type BulkUpsertItemParams struct {
	ItemID    int32
	HistoryID int32
	Deleted   bool
}

func (q *Queries) BulkUpsertItems(ctx context.Context, arg []BulkUpsertItemParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkUpsertItemsStart
	var params []any
	for _, data := range arg {
		query += bulkUpsertItemsValue
		params = append(params,
			data.ItemID,
			data.HistoryID,
			data.Deleted)
	}
	query = strings.TrimRight(query, ",")

	query += bulkUpsertItemsDuplicate

	return q.db.ExecContext(ctx, query, params...)
}
