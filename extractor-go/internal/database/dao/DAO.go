package dao

import (
	"context"
	"database/sql"
)

type IDAO interface {
	Querier

	// i18n.sql.extra.go
	BulkInsertI18nHistory(ctx context.Context, arg []BulkInsertI18nHistoryParams) (sql.Result, error)
	BulkUpsertI18ns(ctx context.Context, arg []BulkUpsertI18nParams) (sql.Result, error)

	// item.sql.extra.go
	BulkInsertItemHistory(ctx context.Context, arg []BulkInsertItemHistoryParams) (sql.Result, error)
	BulkUpsertItems(ctx context.Context, arg []BulkUpsertItemParams) (sql.Result, error)

	// maps.sql.extra.go
	BulkInsertMapHistory(ctx context.Context, arg []BulkInsertMapHistoryParams) (sql.Result, error)
	BulkUpsertMaps(ctx context.Context, arg []BulkUpsertMapParams) (sql.Result, error)

	// quests.sql.extra.go
	BulkInsertQuestHistory(ctx context.Context, arg []BulkInsertQuestHistoryParams) (sql.Result, error)
	BulkUpsertQuests(ctx context.Context, arg []BulkUpsertQuestParams) (sql.Result, error)

	// skills.sql.extra.go
	BulkInsertSkillJob(ctx context.Context, arg []BulkInsertSkillJobParams) (sql.Result, error)
	BulkInsertSkillHistory(ctx context.Context, arg []BulkInsertSkillHistoryParams) (sql.Result, error)
	BulkUpsertSkills(ctx context.Context, arg []BulkUpsertSkillParams) (sql.Result, error)
}
