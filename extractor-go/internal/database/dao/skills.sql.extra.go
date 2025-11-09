// Extensions for quests.sql.go
package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const bulkInsertSkillJobStart = `-- name: InsertSkillJob :execresult
INSERT INTO ` + "`" + `skills_jobs` + "`" + ` (
	` + "`" + `constant` + "`" + `, -- 1
	` + "`" + `job_id` + "`" + `, -- 2
	` + "`" + `inherited_job` + "`" + `, -- 3
	` + "`" + `inherited_job2` + "`" + `, -- 4
	` + "`" + `first_update` + "`" + `, -- 5
	` + "`" + `last_update` + "`" + `, -- 6
	` + "`" + `deleted` + "`" + ` -- 7
)
VALUES
`

const bulkInsertSkillJobValueBlock = `(?, ?, ?, ?, ?, ?, ?),`

const bulkUpsertSkillJobDuplicate = `
ON DUPLICATE KEY UPDATE
	` + "`skills_jobs`.`job_id` = VALUES(`skills_jobs`.`job_id`)," + `
	` + "`skills_jobs`.`inherited_job` = VALUES(`skills_jobs`.`inherited_job`)," + `
	` + "`skills_jobs`.`inherited_job2` = VALUES(`skills_jobs`.`inherited_job2`)," + `
	` + "`skills_jobs`.`first_update` = VALUES(`skills_jobs`.`first_update`)," + `
	` + "`skills_jobs`.`last_update` = VALUES(`skills_jobs`.`last_update`)," + `
	` + "`skills_jobs`.`deleted` = VALUES(`skills_jobs`.`deleted`)"

type BulkInsertSkillJobParams struct {
	Constant      string
	JobId         int32
	InheritedJob  sql.NullInt32
	InheritedJob2 sql.NullInt32
	FirstUpdate   string
	LastUpdate    string
	Deleted       bool
}

func (q *Queries) BulkInsertSkillJob(ctx context.Context, arg []BulkInsertSkillJobParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkInsertSkillJobStart
	var params []any
	for _, data := range arg {
		query += bulkInsertSkillJobValueBlock
		params = append(params,
			data.Constant,
			data.JobId,
			data.InheritedJob,
			data.InheritedJob2,
			data.FirstUpdate,
			data.LastUpdate,
			data.Deleted)
	}
	query = strings.TrimRight(query, ",")
	query += bulkUpsertSkillJobDuplicate

	fmt.Println(query)
	return q.db.ExecContext(ctx, query, params...)
}
