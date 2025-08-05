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

func (r *MapRepository) CountChangesInUpdate(tx *sql.Tx, update string) (int, error) {
	queries := r.DB.GetQueries(tx)
	res, err := queries.CountChangedMapsInUpdate(context.Background(), update)
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *MapRepository) sqlRecordToDomain(dbFrom dao.PreviousMapHistoryVw, dbTo dao.MapsHistory, lastUpdate sql.NullString) FromToRecord[domain.Map] {
	var from *domain.Record[domain.Map] = nil
	var to *domain.Record[domain.Map] = nil

	if dbFrom.HistoryID.Valid {
		from = &domain.Record[domain.Map]{
			Update: dbFrom.Update.String,
			Data:   dbFrom.ToDomain(),
		}
	}

	if dbTo.HistoryID != 0 {
		to = &domain.Record[domain.Map]{
			Update: dbTo.Update,
			Data:   dbTo.ToDomain(),
		}
	}

	return FromToRecord[domain.Map]{
		LastUpdate: domain.NullableString(lastUpdate),
		From:       from,
		To:         to,
	}
}

func (r *MapRepository) GetChangesInUpdate(tx *sql.Tx, update string, pagination Pagination) ([]FromToRecord[domain.Map], error) {
	queries := r.DB.GetQueries(tx)
	res, err := queries.GetChangedMaps(context.Background(), dao.GetChangedMapsParams{
		Update: update,
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []FromToRecord[domain.Map]{}, nil
	}

	if err != nil {
		return nil, err
	}

	records := make([]FromToRecord[domain.Map], len(res))
	for idx, qmodel := range res {
		records[idx] = r.sqlRecordToDomain(qmodel.PreviousMapHistoryVw, qmodel.MapsHistory, qmodel.Lastupdate)
	}

	return records, nil
}

func (r *MapRepository) GetMapHistory(tx *sql.Tx, mapId string, pagination Pagination) ([]FromToRecord[domain.Map], error) {
	queries := r.DB.GetQueries(tx)
	res, err := queries.GetMapHistory(context.Background(), dao.GetMapHistoryParams{
		MapID:  mapId,
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []FromToRecord[domain.Map]{}, nil
	}

	if err != nil {
		return nil, err
	}

	records := make([]FromToRecord[domain.Map], len(res))
	for idx, qmodel := range res {
		records[idx] = r.sqlRecordToDomain(qmodel.PreviousMapHistoryVw, qmodel.MapsHistory, sql.NullString{})
	}

	return records, nil
}

func (r *MapRepository) CountMaps(tx *sql.Tx) (int32, error) {
	queries := r.DB.GetQueries(tx)

	res, err := queries.CountMaps(context.Background())
	if err == sql.ErrNoRows {
		return int32(res), nil
	}

	if err != nil {
		return 0, err
	}

	return int32(res), nil
}

func (r *MapRepository) GetMaps(tx *sql.Tx, pagination Pagination) ([]domain.MinMap, error) {
	queries := r.DB.GetQueries(tx)
	res, err := queries.GetMapList(context.Background(), dao.GetMapListParams{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []domain.MinMap{}, nil
	}
	if err != nil {
		return []domain.MinMap{}, nil
	}

	maps := make([]domain.MinMap, len(res))
	for idx, val := range res {
		maps[idx] = domain.MinMap{
			MapID:      val.MapID,
			LastUpdate: dao.ToNullableString(val.Lastupdate),
			Name:       domain.NullableString(val.Name),
		}
	}
	return maps, nil
}
