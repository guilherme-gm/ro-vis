package database

import (
	"database/sql"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
)

// IDatabase defines the minimal DB contract used by repositories.
// Implement this interface to provide a mockable database in tests.
type IDatabase interface {
	// GetDAO returns the sqlc-generated DAO, optionally bound to a transaction.
	GetDAO(tx *sql.Tx) dao.IDAO
	// BeginTx starts a new transaction.
	BeginTx() (*sql.Tx, error)
}
