package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
)

type LoaderControllerRepository struct {
	DB *database.Database
}

// NewLoaderControllerRepository creates a new LoaderControllerRepository instance
func NewLoaderControllerRepository(db *database.Database) *LoaderControllerRepository {
	return &LoaderControllerRepository{
		DB: db,
	}
}

func (r *LoaderControllerRepository) GetLatestUpdate(tx *sql.Tx, name string) (time.Time, error) {
	queries := r.DB.GetQueries(tx)
	res, err := queries.GetLatestUpdate(context.Background(), name)
	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}

	if err != nil {
		return time.Time{}, err
	}

	return res, nil
}

func (r *LoaderControllerRepository) SetLatestPatch(tx *sql.Tx, name string, date time.Time) error {
	queries := r.DB.GetQueries(tx)
	err := queries.UpsertLatestUpdate(context.Background(), dao.UpsertLatestUpdateParams{
		Name:           name,
		LastUpdateDate: date,
	})
	if err != nil {
		return err
	}

	return nil
}
