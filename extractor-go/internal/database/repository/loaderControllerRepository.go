package repository

import (
	"context"
	"database/sql"

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

func (r *LoaderControllerRepository) GetLatestPatch(name string) (int32, error) {
	res, err := r.queries.GetLatestPatch(context.Background(), name)
	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return res, nil
}

func (r *LoaderControllerRepository) SetLatestPatch(name string, patchId int32, patchName string) error {
	err := r.queries.UpsertLatestPatch(context.Background(), dao.UpsertLatestPatchParams{
		Name:            name,
		LatestPatchID:   patchId,
		LatestPatchName: patchName,
	})
	if err != nil {
		return err
	}

	return nil
}
