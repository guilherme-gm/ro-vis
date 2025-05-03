package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type PatchRepository struct {
	queries *dao.Queries
}

func newPatchRepository(queries *dao.Queries) *PatchRepository {
	return &PatchRepository{queries: queries}
}

func GetPatchRepository() *PatchRepository {
	if repositoriesCache.patchRepository == nil {
		repositoriesCache.patchRepository = newPatchRepository(database.GetQueries())
	}

	return repositoriesCache.patchRepository
}

func (r *PatchRepository) ListPatches() (*[]domain.Patch, error) {
	res, err := r.queries.ListPatches(context.Background())
	if err == sql.ErrNoRows {
		return &[]domain.Patch{}, nil
	}

	if err != nil {
		return nil, err
	}

	patches := make([]domain.Patch, len(res))
	for idx, pmodel := range res {
		patches[idx] = pmodel.ToDomain()
	}

	return &patches, nil
}

func (r *PatchRepository) InsertPatch(patch *domain.Patch) error {
	jsonMsg, err := json.Marshal(patch.Files)
	if err != nil {
		return err
	}

	return r.queries.InsertPatch(context.Background(), dao.InsertPatchParams{
		Name:  patch.Name,
		Date:  patch.Date,
		Files: jsonMsg,
	})
}
