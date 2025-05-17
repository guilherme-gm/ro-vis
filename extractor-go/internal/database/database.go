package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
)

var db *sql.DB = nil

func GetDB() *sql.DB {
	if db == nil {
		dbConn, err := sql.Open("mysql", conf.Config.DbUrl)
		if err != nil {
			panic(err)
		}

		db = dbConn
	}

	return db
}

func GetQueries(tx *sql.Tx) *dao.Queries {
	conn := GetDB()
	queries := dao.New(conn)
	if tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries
}

func BeginTx() (*sql.Tx, error) {
	return GetDB().Begin()
}
