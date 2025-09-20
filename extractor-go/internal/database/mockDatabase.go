package database

import (
	"database/sql"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/stretchr/testify/mock"
)

type mockDatabase struct {
	Dao *dao.MockIDAO
}

func NewMockDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDatabase {
	mock := &mockDatabase{}
	mock.Dao = dao.NewMockIDAO(t)

	return mock
}

func (m *mockDatabase) GetDAO(tx *sql.Tx) dao.IDAO {
	return m.Dao
}

func (m *mockDatabase) BeginTx() (*sql.Tx, error) {
	return nil, nil
}
