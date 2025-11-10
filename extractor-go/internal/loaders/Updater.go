package loaders

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

/**
 * Generic updater to track changes to a given data type.
 *
 * It is used to generically track changes in loaders so they can later apply it to the database.
 */

type Updater[K comparable, T domain.Identifiable[K], P domain.IdentifiablePointer[K, T]] struct {
	CurrentValues  map[K]P
	DirtyValues    map[K]P
	ValuesToInsert []P
	ValuesToUpdate []P
}

func NewUpdater[K comparable, T domain.Identifiable[K], P domain.IdentifiablePointer[K, T]](currentValues []T) *Updater[K, T, P] {
	currentValuesHash := make(map[K]P)
	for _, m := range currentValues {
		currentValuesHash[m.GetId()] = &m
	}

	return &Updater[K, T, P]{
		DirtyValues:   make(map[K]P),
		CurrentValues: currentValuesHash,
	}
}

func (u *Updater[K, T, P]) GetForRead(key K) (T, bool) {
	if m, ok := u.DirtyValues[key]; ok {
		return *m, true
	}

	if m, ok := u.CurrentValues[key]; ok {
		return *m, true
	}

	var empty T
	return empty, false
}

func (u *Updater[K, T, P]) GetForEdit(key K) P {
	if m, ok := u.DirtyValues[key]; ok {
		return m
	}

	if m, ok := u.CurrentValues[key]; ok {
		newMap := *m
		u.ValuesToUpdate = append(u.ValuesToUpdate, &newMap)
		u.DirtyValues[key] = &newMap
		return &newMap
	}

	var newMapData T
	newMap := P(&newMapData)
	newMap.SetId(key)
	u.ValuesToInsert = append(u.ValuesToInsert, newMap)
	u.DirtyValues[key] = newMap
	return newMap
}

func (u *Updater[K, T, P]) ForEach(fn func(K, T)) {
	readKeys := make(map[K]bool)

	// Dirty contains modified/new values
	for _, m := range u.DirtyValues {
		if _, ok := readKeys[(*m).GetId()]; ok {
			continue
		}

		readKeys[(*m).GetId()] = true
		fn((*m).GetId(), *m)
	}

	// Current contains existing values, before any changes (may be duplicated with dirty)
	for _, m := range u.CurrentValues {
		if _, ok := readKeys[(*m).GetId()]; ok {
			continue
		}
		readKeys[(*m).GetId()] = true
		fn((*m).GetId(), *m)
	}
}
