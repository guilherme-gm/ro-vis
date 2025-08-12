package database

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
)

type Database struct {
	Connection *sql.DB
}

func NewDatabase(databaseName string) *Database {
	dbUrl := strings.Replace(conf.Config.DbUrl, "@DATABASE@", databaseName, 1)
	dbConn, err := sql.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}

	return &Database{Connection: dbConn}
}

func (db *Database) GetQueries(tx *sql.Tx) *dao.Queries {
	conn := db.Connection
	queries := dao.New(conn)
	if tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries
}

func (db *Database) BeginTx() (*sql.Tx, error) {
	return db.Connection.Begin()
}
