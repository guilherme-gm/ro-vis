package mapData

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

type MapLoader struct {
	repository *repository.MapRepository
	server     *server.Server
}

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

// GetRelevantFiles returns a list of all files that are relevant to this loader's parsers.
func (l *MapLoader) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		mapNameTableRegex,
		mp3NameTableRegex,
		naviMapRegex,
		naviNpcRegex,
		naviMobRegex,
		naviLinkRegex,
	}
}

func NewMapLoader(server *server.Server) *MapLoader {
	return &MapLoader{
		repository: server.Repositories.MapRepository,
		server:     server,
	}
}

func (l *MapLoader) LoadPatch(tx *sql.Tx, basePath string, update domain.Update) {
	if !update.HasChangedAnyFiles(l.GetRelevantFiles()) {
		fmt.Println("Skipped - No meaningful file")
		return
	}

	fmt.Println("> Decoding...")

	fmt.Println("> Fetching current list...")
	currentMaps, err := l.repository.GetCurrentMaps(tx)
	if err != nil {
		panic(err)
	}

	updater := newMapUpdater(*currentMaps)
	l.loadMapNames(basePath, update, updater)
	l.loadMp3Names(basePath, update, updater)
	l.loadSpecialCodes(basePath, update, updater)
	l.loadNpcs(basePath, update, updater)
	l.loadSpawns(basePath, update, updater)
	l.loadWarps(basePath, update, updater)

	// Saving
	fmt.Println("> Saving...")
	fmt.Printf("> Saving new records... (%d records to save)\n", len(updater.mapsToInsert))
	err = l.repository.AddMapsToHistory(tx, update.Name(), updater.mapsToInsert)
	if err != nil {
		fmt.Printf("Error saving new records: %v\n", err)
		panic(err)
	}

	fmt.Printf("> Updating records... (%d records to update)\n", len(updater.mapsToUpdate))
	err = l.repository.AddMapsToHistory(tx, update.Name(), updater.mapsToUpdate)
	if err != nil {
		panic(err)
	}
}

func (l *MapLoader) loadMapNames(basePath string, update domain.Update, updater *mapUpdater) {
	change, err := update.GetChangeForFile(mapNameTable)
	if err != nil {
		if errors.Is(err, domain.NewNotFoundError("")) {
			return
		}

		panic(err)
	}

	mapNames, err := ParseMapValueTable(basePath + "/" + change.Patch + "/" + change.File)
	if err != nil {
		panic(err)
	}

	for mapId, name := range mapNames {
		mapId = strings.TrimSuffix(mapId, ".rsw")

		exitingMap := updater.getForRead(mapId)
		nameObj := domain.NullableString{String: name, Valid: true}
		if exitingMap.Id == "" || exitingMap.Name != nameObj {
			updater.getForEdit(mapId).Name = nameObj
		}
	}
}

func (l *MapLoader) loadMp3Names(basePath string, update domain.Update, updater *mapUpdater) {
	change, err := update.GetChangeForFile(mp3NameTable)
	if err != nil {
		if errors.Is(err, domain.NewNotFoundError("")) {
			return
		}

		panic(err)
	}

	mp3Names, err := ParseMapValueTable(basePath + "/" + change.Patch + "/" + change.File)
	if err != nil {
		panic(err)
	}

	for mapId, name := range mp3Names {
		mapId = strings.TrimSuffix(mapId, ".rsw")
		name = strings.ReplaceAll(name, `\\`, "/")

		existingMap := updater.getForRead(mapId)
		nameObj := domain.NullableString{String: name, Valid: true}
		if existingMap.Id == "" || existingMap.Mp3Name != nameObj {
			updater.getForEdit(mapId).Mp3Name = nameObj
		}
	}
}

func (l *MapLoader) loadSpecialCodes(basePath string, update domain.Update, updater *mapUpdater) {
	change, err := update.GetChangeForFile(naviMap)
	if err != nil {
		if errors.Is(err, domain.NewNotFoundError("")) {
			return
		}

		panic(err)
	}

	naviMapParser := NewNaviMapV1Parser(l.server)
	naviMaps := naviMapParser.Parse(basePath, &change)

	for _, naviMap := range naviMaps {
		exitingMap := updater.getForRead(naviMap.MapId)
		specialCodeObj := domain.NullableInt32{Int32: int32(naviMap.SpecialCode), Valid: true}
		if exitingMap.Id == "" || exitingMap.SpecialCode != specialCodeObj {
			updater.getForEdit(naviMap.MapId).SpecialCode = specialCodeObj
		}
	}
}

func (l *MapLoader) loadNpcs(basePath string, update domain.Update, updater *mapUpdater) {
	change, err := update.GetChangeForFile(naviNpc)
	if err != nil {
		if errors.Is(err, domain.NewNotFoundError("")) {
			return
		}

		panic(err)
	}

	naviNpcParser := NewNaviNpcV1Parser(l.server)
	naviNpcs := naviNpcParser.Parse(basePath, &change)
	mapToNpcs := make(map[string][]domain.MapNpc)

	for _, naviNpc := range naviNpcs {
		mapToNpcs[naviNpc.MapId] = append(mapToNpcs[naviNpc.MapId], naviNpc.ToDomain())
	}

	for _, existingMap := range updater.currentMaps {
		newMapNpcs, ok := mapToNpcs[existingMap.Id]
		if !ok {
			// This map does not have NPCs.

			if len(existingMap.Npcs) > 0 {
				// But we currently have NPCs for it. Remove them.
				updater.getForEdit(existingMap.Id).Npcs = nil
			}

			// otherwise, nothing to do here.
			continue
		}

		if !areArraysEqual(existingMap.Npcs, newMapNpcs) {
			updater.getForEdit(existingMap.Id).Npcs = newMapNpcs
		}

		// This map has been checked, we can remove from the list
		delete(mapToNpcs, existingMap.Id)
	}

	// Any remaining map in the list is a new map
	for mapId, npcs := range mapToNpcs {
		updater.getForEdit(mapId).Npcs = npcs
	}
}

func (l *MapLoader) loadSpawns(basePath string, update domain.Update, updater *mapUpdater) {
	change, err := update.GetChangeForFile(naviMob)
	if err != nil {
		if errors.Is(err, domain.NewNotFoundError("")) {
			return
		}

		panic(err)
	}

	naviMobParser := NewNaviMobV1Parser(l.server)
	naviMobs := naviMobParser.Parse(basePath, &change)
	mapToSpawns := make(map[string][]domain.MapSpawn)

	for _, naviMob := range naviMobs {
		mapToSpawns[naviMob.MapId] = append(mapToSpawns[naviMob.MapId], naviMob.ToDomain())
	}

	for _, existingMap := range updater.currentMaps {
		newMapSpawns, ok := mapToSpawns[existingMap.Id]
		if !ok {
			// This map does not have spawns.

			if len(existingMap.Spawns) > 0 {
				// But we currently have spawns for it. Remove them.
				updater.getForEdit(existingMap.Id).Spawns = nil
			}

			// otherwise, nothing to do here.
			continue
		}

		if !areArraysEqual(existingMap.Spawns, newMapSpawns) {
			updater.getForEdit(existingMap.Id).Spawns = newMapSpawns
		}

		// This map has been checked, we can remove from the list
		delete(mapToSpawns, existingMap.Id)
	}

	// Any remaining map in the list is a new map
	for mapId, spawns := range mapToSpawns {
		updater.getForEdit(mapId).Spawns = spawns
	}
}

func (l *MapLoader) loadWarps(basePath string, update domain.Update, updater *mapUpdater) {
	change, err := update.GetChangeForFile(naviLink)
	if err != nil {
		if errors.Is(err, domain.NewNotFoundError("")) {
			return
		}

		panic(err)
	}

	naviLinkParser := NewNaviLinkV1Parser(l.server)
	naviLinks := naviLinkParser.Parse(basePath, &change)
	mapToWarps := make(map[string][]domain.MapWarp)

	for _, naviLink := range naviLinks {
		mapToWarps[naviLink.MapId] = append(mapToWarps[naviLink.MapId], naviLink.ToDomain())
	}

	for _, existingMap := range updater.currentMaps {
		newMapWarps, ok := mapToWarps[existingMap.Id]
		if !ok {
			// This map does not have warps.

			if len(existingMap.Warps) > 0 {
				// But we currently have warps for it. Remove them.
				updater.getForEdit(existingMap.Id).Warps = nil
			}

			// otherwise, nothing to do here.
			continue
		}

		if !areArraysEqual(existingMap.Warps, newMapWarps) {
			updater.getForEdit(existingMap.Id).Warps = newMapWarps
		}

		// This map has been checked, we can remove from the list
		delete(mapToWarps, existingMap.Id)
	}

	// Any remaining map in the list is a new map
	for mapId, warps := range mapToWarps {
		updater.getForEdit(mapId).Warps = warps
	}
}

func areArraysEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	// This is not the safest thing to do here, but I am ok with it for now
	// If for some reason the files changes order, this will hit as different
	// but I think it is fine for now.
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func (l *MapLoader) Name() string {
	return "mapData"
}
