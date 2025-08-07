package mapData

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type mapUpdater struct {
	currentMaps  map[string]*domain.Map
	dirtyMaps    map[string]*domain.Map
	mapsToInsert []*domain.Map
	mapsToUpdate []*domain.Map
}

func newMapUpdater(currentMaps []domain.Map) *mapUpdater {
	currentMapHash := make(map[string]*domain.Map)
	for _, m := range currentMaps {
		currentMapHash[m.Id] = &m
	}

	return &mapUpdater{
		dirtyMaps:   make(map[string]*domain.Map),
		currentMaps: currentMapHash,
	}
}

func (u *mapUpdater) getForRead(mapId string) domain.Map {
	if m, ok := u.dirtyMaps[mapId]; ok {
		return *m
	}

	if m, ok := u.currentMaps[mapId]; ok {
		return *m
	}

	return domain.Map{}
}

func (u *mapUpdater) getForEdit(mapId string) *domain.Map {
	if m, ok := u.dirtyMaps[mapId]; ok {
		return m
	}

	if m, ok := u.currentMaps[mapId]; ok {
		newMap := *m
		newMap.PreviousHistoryID = m.HistoryID
		u.mapsToUpdate = append(u.mapsToUpdate, &newMap)
		u.dirtyMaps[mapId] = &newMap
		return &newMap
	}

	newMap := domain.Map{
		Id: mapId,
	}
	u.mapsToInsert = append(u.mapsToInsert, &newMap)
	u.dirtyMaps[mapId] = &newMap
	return &newMap
}
