// Extensions for quests.sql.go
package dao

import (
	"context"
	"database/sql"
	"strings"
)

// ============== Bulk Insert Skill Job ===================

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

	return q.db.ExecContext(ctx, query, params...)
}

// =================== Bulk Insert Skill History ==================

const bulkInsertSkillHistoryStart = `-- name: InsertSkillHistory :execresult
INSERT INTO ` + "`" + `skills_history` + "`" + ` (
	` + "`" + `previous_history_id` + "`" + `, -- 1
	` + "`" + `skill_id` + "`" + `, -- 2
	` + "`" + `file_version` + "`" + `, -- 3
	` + "`" + `update` + "`" + `, -- 4
	` + "`" + `constant` + "`" + `, -- 5
	` + "`" + `name` + "`" + `, -- 6
	` + "`" + `description` + "`" + `, -- 7
	` + "`" + `max_level` + "`" + `, -- 8
	` + "`" + `sp_cost` + "`" + `, -- 9
	` + "`" + `ap_cost` + "`" + `, -- 10
	` + "`" + `can_select_level` + "`" + `, -- 11
	` + "`" + `attack_range` + "`" + `, -- 12
	` + "`" + `required_skills` + "`" + `, -- 13
	` + "`" + `job_required_skills` + "`" + `, -- 14
	` + "`" + `skill_scale` + "`" + ` -- 15
)
VALUES
`

const bulkInsertSkillHistoryValueBlock = `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),`

type BulkInsertSkillHistoryParams struct {
	PreviousHistoryID sql.NullInt32
	SkillId           int32
	FileVersion       int32
	Update            string
	Constant          sql.NullString
	Name              sql.NullString
	Description       sql.NullString
	MaxLevel          sql.NullInt32
	SpCost            sql.NullString
	ApCost            sql.NullString
	CanSelectLevel    sql.NullBool
	AttackRange       sql.NullString
	RequiredSkills    sql.NullString
	JobRequiredSkills sql.NullString
	SkillScale        sql.NullString
}

func (q *Queries) BulkInsertSkillHistory(ctx context.Context, arg []BulkInsertSkillHistoryParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkInsertSkillHistoryStart
	var params []any
	for _, data := range arg {
		query += bulkInsertSkillHistoryValueBlock
		params = append(params,
			data.PreviousHistoryID,
			data.SkillId,
			data.FileVersion,
			data.Update,
			data.Constant,
			data.Name,
			data.Description,
			data.MaxLevel,
			data.SpCost,
			data.ApCost,
			data.CanSelectLevel,
			data.AttackRange,
			data.RequiredSkills,
			data.JobRequiredSkills,
			data.SkillScale)
	}
	query = strings.TrimRight(query, ",")
	return q.db.ExecContext(ctx, query, params...)
}

// =================== Bulk Upsert Skill ==================

const bulkUpsertSkillHistoryStart = `-- name: BulkUpsertSkillHistory :execresult
INSERT INTO ` + "`skills` (`skill_id`,`latest_history_id`,`deleted`)" + `
VALUES
`

const bulkUpsertSkillHistoryValue = `(?, ?, ?),`

const bulkUpsertSkillHistoryDuplicate = `
ON DUPLICATE KEY UPDATE
	` + "`skills`.`latest_history_id` = VALUES(`skills`.`latest_history_id`)," + `
	` + "`skills`.`deleted` = VALUES(`skills`.`deleted`)"

type BulkUpsertSkillParams struct {
	SkillId   int32
	HistoryID int32
	Deleted   bool
}

func (q *Queries) BulkUpsertSkills(ctx context.Context, arg []BulkUpsertSkillParams) (sql.Result, error) {
	if len(arg) == 0 {
		return nil, sql.ErrNoRows
	}

	query := bulkUpsertSkillHistoryStart
	var params []any
	for _, data := range arg {
		query += bulkUpsertSkillHistoryValue
		params = append(params,
			data.SkillId,
			data.HistoryID,
			data.Deleted)
	}
	query = strings.TrimRight(query, ",")

	query += bulkUpsertSkillHistoryDuplicate

	return q.db.ExecContext(ctx, query, params...)
}
