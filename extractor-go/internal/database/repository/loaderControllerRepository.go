package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
)

type LoaderControllerRepository struct {
	queries *dao.Queries
}

func newLoaderControllerRepository(queries *dao.Queries) *LoaderControllerRepository {
	return &LoaderControllerRepository{queries: queries}
}

func GetLoaderControllerRepository() *LoaderControllerRepository {
	if repositoriesCache.LoaderControllerRepository == nil {
		repositoriesCache.LoaderControllerRepository = newLoaderControllerRepository(database.GetQueries())
	}

	return repositoriesCache.LoaderControllerRepository
}

func (r *LoaderControllerRepository) GetLatestUpdate(name string) (time.Time, error) {
	res, err := r.queries.GetLatestUpdate(context.Background(), name)
	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}

	if err != nil {
		return time.Time{}, err
	}

	return res, nil
}

func (r *LoaderControllerRepository) SetLatestPatch(name string, date time.Time) error {
	err := r.queries.UpsertLatestUpdate(context.Background(), dao.UpsertLatestUpdateParams{
		Name:           name,
		LastUpdateDate: date,
	})
	if err != nil {
		return err
	}

	return nil
}
