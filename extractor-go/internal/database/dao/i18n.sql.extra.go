// Extensions for quests.sql.go
package dao

import (
	"context"
	"database/sql"
	"strings"
)

const bulkInsertI18nHistoryStart = `-- name: InsertI18nHistory :execresult
INSERT INTO ` + "`" + `i18n_history` + "`" + ` (
	` + "`" + `previous_history_id` + "`" + `, -- 1
	` + "`" + `i18n_id` + "`" + `, -- 2
	` + "`" + `file_version` + "`" + `, -- 3
	` + "`" + `update` + "`" + `, -- 4
	` + "`" + `container_file` + "`" + `, -- 5
	` + "`" + `en_text` + "`" + `, -- 6
	` + "`" + `pt_br_text` + "`" + `, -- 7
	` + "`" + `active` + "`" + ` -- 8
)
VALUES
`

const bulkInsertI18nHistoryValueBlock = `(?, ?, ?, ?, ?, ?, ?, ?),`

type BulkInsertI18nHistoryParams struct {
	PreviousHistoryID sql.NullInt64
	I18nID            string
	FileVersion       int32
	Update            string
	ContainerFile     string
	EnText            string
	PtBrText          string
	Active            bool
}

func (q *Queries) BulkInsertI18nHistory(ctx context.Context, arg []BulkInsertI18nHistoryParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkInsertI18nHistoryStart
	var params []any
	for _, data := range arg {
		query += bulkInsertI18nHistoryValueBlock
		params = append(params,
			data.PreviousHistoryID,
			data.I18nID,
			data.FileVersion,
			data.Update,
			data.ContainerFile,
			data.EnText,
			data.PtBrText,
			data.Active)
	}
	query = strings.TrimRight(query, ",")
	return q.db.ExecContext(ctx, query, params...)
}

const bulkUpsertI18nsStart = `-- name: BulkUpsertI18nHistory :execresult
INSERT INTO ` + "`i18ns` (`i18n_id`,`latest_history_id`,`deleted`)" + `
VALUES
`

const bulkUpsertI18nsValue = `(?, ?, ?),`

const bulkUpsertI18nsDuplicate = `
ON DUPLICATE KEY UPDATE
	` + "`i18ns`.`latest_history_id` = VALUES(`i18ns`.`latest_history_id`)," + `
	` + "`i18ns`.`deleted` = VALUES(`i18ns`.`deleted`)"

type BulkUpsertI18nParams struct {
	I18nID    string
	HistoryID int64
	Deleted   bool
}

func (q *Queries) BulkUpsertI18ns(ctx context.Context, arg []BulkUpsertI18nParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkUpsertI18nsStart
	var params []any
	for _, data := range arg {
		query += bulkUpsertI18nsValue
		params = append(params,
			data.I18nID,
			data.HistoryID,
			data.Deleted)
	}
	query = strings.TrimRight(query, ",")

	query += bulkUpsertI18nsDuplicate

	return q.db.ExecContext(ctx, query, params...)
}
