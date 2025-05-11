package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

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

func (r *PatchRepository) ListPatches() ([]domain.Patch, error) {
	res, err := r.queries.ListPatches(context.Background())
	if err == sql.ErrNoRows {
		return []domain.Patch{}, nil
	}

	if err != nil {
		return nil, err
	}

	patches := make([]domain.Patch, len(res))
	for idx, pmodel := range res {
		patches[idx] = pmodel.ToDomain()
	}

	return patches, nil
}

func (r *PatchRepository) ListUpdates() ([]domain.Update, error) {
	patches, err := r.ListPatches()
	if err != nil || len(patches) == 0 {
		return []domain.Update{}, err
	}

	updates := []domain.Update{}

	var time time.Time
	var update *domain.Update
	fileMap := make(map[string]*domain.UpdateChange)
	for _, patch := range patches {
		if patch.Date != time {
			time = patch.Date

			updates = append(updates, domain.Update{
				Date: time,
			})
			update = &updates[len(updates)-1]

			for k := range fileMap {
				delete(fileMap, k)
			}
		}

		for _, file := range patch.Files {
			ptr, ok := fileMap[file]
			if ok {
				ptr.Patch = patch.Name
				ptr.File = file
			} else {
				change := domain.UpdateChange{
					Patch: patch.Name,
					File:  file,
				}

				update.Changes = append(update.Changes, change)
				fileMap[file] = &update.Changes[len(update.Changes)-1]
			}
		}
	}

	return updates, nil
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
