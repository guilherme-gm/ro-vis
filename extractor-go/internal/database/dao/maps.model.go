package dao

import (
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func (q *GetCurrentMapsRow) ToDomain() domain.Map {
	var npcs []domain.MapNpc
	var warps []domain.MapWarp
	var spawns []domain.MapSpawn

	json.Unmarshal(q.Npcs, &npcs)
	json.Unmarshal(q.Warps, &warps)
	json.Unmarshal(q.Spawns, &spawns)

	return domain.Map{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         ToNullableInt32(q.HistoryID),
		Id:                q.MapID,
		FileVersion:       q.FileVersion,
		Name:              domain.NullableString(q.Name),
		SpecialCode:       domain.NullableInt32(q.SpecialCode),
		Mp3Name:           domain.NullableString(q.Mp3Name),
		Npcs:              npcs,
		Warps:             warps,
		Spawns:            spawns,
	}
}

func (q *MapsHistory) ToDomain() domain.Map {
	var npcs []domain.MapNpc
	var warps []domain.MapWarp
	var spawns []domain.MapSpawn

	json.Unmarshal(q.Npcs, &npcs)
	json.Unmarshal(q.Warps, &warps)
	json.Unmarshal(q.Spawns, &spawns)

	return domain.Map{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         ToNullableInt32(q.HistoryID),
		Id:                q.MapID,
		FileVersion:       q.FileVersion,
		Name:              domain.NullableString(q.Name),
		SpecialCode:       domain.NullableInt32(q.SpecialCode),
		Mp3Name:           domain.NullableString(q.Mp3Name),
		Npcs:              npcs,
		Warps:             warps,
		Spawns:            spawns,
	}
}

func (q *PreviousMapHistoryVw) ToDomain() domain.Map {
	var npcs []domain.MapNpc
	var warps []domain.MapWarp
	var spawns []domain.MapSpawn

	json.Unmarshal(q.Npcs, &npcs)
	json.Unmarshal(q.Warps, &warps)
	json.Unmarshal(q.Spawns, &spawns)

	return domain.Map{
		PreviousHistoryID: ToNullableInt32(q.PreviousHistoryID.Int32),
		HistoryID:         ToNullableInt32(q.HistoryID.Int32),
		Id:                q.MapID.String,
		Name:              domain.NullableString(q.Name),
		SpecialCode:       domain.NullableInt32(q.SpecialCode),
		Mp3Name:           domain.NullableString(q.Mp3Name),
		Npcs:              npcs,
		Warps:             warps,
		Spawns:            spawns,
	}
}
