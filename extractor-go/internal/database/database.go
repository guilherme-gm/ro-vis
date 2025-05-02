package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
)

var queries *dao.Queries = nil

func GetQueries() *dao.Queries {
	if queries != nil {
		return queries
	}

	db, err := sql.Open("mysql", conf.Config.DbUrl)
	if err != nil {
		panic(err)
	}

	queries = dao.New(db)
	return queries
}
