// Extensions for maps.sql.go
package dao

import (
	"context"
	"database/sql"
	"strings"
)

const bulkInsertMapHistoryStart = `-- name: InsertMapHistory :execresult
INSERT INTO ` + "`" + `maps_history` + "`" + ` (
	` + "`" + `previous_history_id` + "`" + `, -- 1
	` + "`" + `map_id` + "`" + `, -- 2
	` + "`" + `file_version` + "`" + `, -- 3
	` + "`" + `update` + "`" + `, -- 4
	` + "`" + `name` + "`" + `, -- 5
	` + "`" + `special_code` + "`" + `, -- 6
	` + "`" + `mp3_name` + "`" + `, -- 7
	` + "`" + `npcs` + "`" + `, -- 8
	` + "`" + `warps` + "`" + `, -- 9
	` + "`" + `spawns` + "`" + ` -- 10
)
VALUES
`

const bulkInsertMapHistoryValueBlock = `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?),`

type BulkInsertMapHistoryParams struct {
	PreviousHistoryID sql.NullInt32
	MapId             string
	FileVersion       int32
	Update            string
	Name              sql.NullString
	SpecialCode       sql.NullInt32
	MP3Name           sql.NullString
	Npcs              sql.NullString
	Warps             sql.NullString
	Spawns            sql.NullString
}

func (q *Queries) BulkInsertMapHistory(ctx context.Context, arg []BulkInsertMapHistoryParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkInsertMapHistoryStart
	var params []any
	for _, data := range arg {
		query += bulkInsertMapHistoryValueBlock
		params = append(params,
			data.PreviousHistoryID,
			data.MapId,
			data.FileVersion,
			data.Update,
			data.Name,
			data.SpecialCode,
			data.MP3Name,
			data.Npcs,
			data.Warps,
			data.Spawns)
	}
	query = strings.TrimRight(query, ",")
	return q.db.ExecContext(ctx, query, params...)
}

const bulkUpsertMapHistoryStart = `-- name: BulkUpsertMapHistory :execresult
INSERT INTO ` + "`maps` (`map_id`,`latest_history_id`,`deleted`)" + `
VALUES
`

const bulkUpsertMapHistoryValue = `(?, ?, ?),`

const bulkUpsertMapHistoryDuplicate = `
ON DUPLICATE KEY UPDATE
	` + "`maps`.`latest_history_id` = VALUES(`maps`.`latest_history_id`)," + `
	` + "`maps`.`deleted` = VALUES(`maps`.`deleted`)"

type BulkUpsertMapParams struct {
	MapId     string
	HistoryID int32
	Deleted   bool
}

func (q *Queries) BulkUpsertMaps(ctx context.Context, arg []BulkUpsertMapParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkUpsertMapHistoryStart
	var params []any
	for _, data := range arg {
		query += bulkUpsertMapHistoryValue
		params = append(params,
			data.MapId,
			data.HistoryID,
			data.Deleted)
	}
	query = strings.TrimRight(query, ",")

	query += bulkUpsertMapHistoryDuplicate

	return q.db.ExecContext(ctx, query, params...)
}
