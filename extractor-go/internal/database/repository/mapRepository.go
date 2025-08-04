package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type MapRepository struct {
	DB *database.Database
}

// NewMapRepository creates a new MapRepository instance
func NewMapRepository(db *database.Database) *MapRepository {
	return &MapRepository{
		DB: db,
	}
}

func (r *MapRepository) GetCurrentMaps(tx *sql.Tx) (*[]domain.Map, error) {
	queries := r.DB.GetQueries(tx)
	res, err := queries.GetCurrentMaps(context.Background())
	if err == sql.ErrNoRows {
		return &[]domain.Map{}, nil
	}

	if err != nil {
		return nil, err
	}

	maps := make([]domain.Map, len(res))
	for idx, qmodel := range res {
		maps[idx] = qmodel.ToDomain()
	}

	return &maps, nil
}

func (r *MapRepository) addMapsToHistory_sub(tx *sql.Tx, update string, newMaps []*domain.Map) error {
	queries := r.DB.GetQueries(tx)
	bulkParams := []dao.BulkInsertMapHistoryParams{}
	updatedIdMap := make(map[string]bool, len(newMaps))
	for _, it := range newMaps {
		updatedIdMap[it.Id] = true
		var npcsJson string
		if len(it.Npcs) > 0 {
			jsonBytes, _ := json.Marshal(it.Npcs)
			npcsJson = string(jsonBytes)
		}
		var warpsJson string
		if len(it.Warps) > 0 {
			jsonBytes, _ := json.Marshal(it.Warps)
			warpsJson = string(jsonBytes)
		}
		var spawnsJson string
		if len(it.Spawns) > 0 {
			jsonBytes, _ := json.Marshal(it.Spawns)
			spawnsJson = string(jsonBytes)
		}
		bulkParams = append(bulkParams, dao.BulkInsertMapHistoryParams{
			PreviousHistoryID: sql.NullInt32(it.PreviousHistoryID),
			MapId:             it.Id,
			FileVersion:       it.FileVersion,
			Update:            update,
			Name:              sql.NullString(it.Name),
			SpecialCode:       sql.NullInt32(it.SpecialCode),
			MP3Name:           sql.NullString(it.Mp3Name),
			Npcs:              sql.NullString{String: npcsJson, Valid: len(it.Npcs) > 0},
			Warps:             sql.NullString{String: warpsJson, Valid: len(it.Warps) > 0},
			Spawns:            sql.NullString{String: spawnsJson, Valid: len(it.Spawns) > 0},
		})
	}

	_, err := queries.BulkInsertMapHistory(context.Background(), bulkParams)
	if err != nil {
		return err
	}

	res, err := queries.GetMapsIdsInUpdate(context.Background(), update)
	if err != nil {
		return err
	}

	upsertParams := []dao.BulkUpsertMapParams{}
	for _, id := range res {
		if _, ok := updatedIdMap[id.MapID]; !ok {
			continue
		}

		upsertParams = append(upsertParams, dao.BulkUpsertMapParams{
			MapId:     id.MapID,
			HistoryID: id.HistoryID,
			Deleted:   false,
		})
	}

	_, err = queries.BulkUpsertMaps(context.Background(), upsertParams)
	if err != nil {
		return err
	}

	return err
}

func (r *MapRepository) AddMapsToHistory(tx *sql.Tx, update string, newMaps []*domain.Map) error {
	if len(newMaps) == 0 {
		return nil
	}

	steps := (len(newMaps) / 500)

	i := 0
	for ; i < steps; i++ {
		slice := newMaps[i*500 : (i+1)*500]
		if err := r.addMapsToHistory_sub(tx, update, slice); err != nil {
			return err
		}
	}

	slice := newMaps[i*500:]
	if err := r.addMapsToHistory_sub(tx, update, slice); err != nil {
		return err
	}

	return nil
}
