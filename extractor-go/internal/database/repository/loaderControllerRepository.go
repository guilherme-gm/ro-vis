package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
)

type LoaderControllerRepository struct{}

func newLoaderControllerRepository() *LoaderControllerRepository {
	return &LoaderControllerRepository{}
}

func GetLoaderControllerRepository() *LoaderControllerRepository {
	if repositoriesCache.LoaderControllerRepository == nil {
		repositoriesCache.LoaderControllerRepository = newLoaderControllerRepository()
	}

	return repositoriesCache.LoaderControllerRepository
}

func (r *LoaderControllerRepository) GetLatestUpdate(tx *sql.Tx, name string) (time.Time, error) {
	queries := database.GetQueries(tx)
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
	queries := database.GetQueries(tx)
	err := queries.UpsertLatestUpdate(context.Background(), dao.UpsertLatestUpdateParams{
		Name:           name,
		LastUpdateDate: date,
	})
	if err != nil {
		return err
	}

	return nil
}
