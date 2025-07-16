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
	DB *database.Database
}

// NewPatchRepository creates a new PatchRepository instance
func NewPatchRepository(db *database.Database) *PatchRepository {
	return &PatchRepository{
		DB: db,
	}
}

func (r *PatchRepository) ListPatches(tx *sql.Tx) ([]domain.Patch, error) {
	queries := r.DB.GetQueries(tx)
	res, err := queries.ListPatches(context.Background())
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

func (r *PatchRepository) listUpdatesPatches(tx *sql.Tx, pagination Pagination) ([]domain.Patch, error) {
	queries := r.DB.GetQueries(tx)
	res, err := queries.ListUpdatesPatches(context.Background(), dao.ListUpdatesPatchesParams{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
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

func (r *PatchRepository) GetUpdateCount(tx *sql.Tx) (int32, error) {
	count, err := r.DB.GetQueries(tx).GetUpdatesCount(context.Background())
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

func (r *PatchRepository) ListUpdates(tx *sql.Tx, pagination Pagination) ([]domain.Update, error) {
	var patches []domain.Patch
	var err error
	if pagination == PaginateAll {
		patches, err = r.ListPatches(tx)
	} else {
		patches, err = r.listUpdatesPatches(tx, pagination)
	}

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

func (r *PatchRepository) InsertPatch(tx *sql.Tx, patch *domain.Patch) error {
	queries := r.DB.GetQueries(tx)
	jsonMsg, err := json.Marshal(patch.Files)
	if err != nil {
		return err
	}

	return queries.InsertPatch(context.Background(), dao.InsertPatchParams{
		Name:  patch.Name,
		Date:  patch.Date,
		Files: jsonMsg,
	})
}
