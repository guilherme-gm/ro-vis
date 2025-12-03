package loaders

import "github.com/guilherme-gm/ro-vis/extractor/internal/domain"

type UpdaterEntry[K comparable] interface {
	GetId() K
	GetHistoryId() domain.NullableInt32
}

type UpdaterEntryPointer[K comparable, T any] interface {
	SetId(id K)
	SetHistoryId(id domain.NullableInt32)
	SetPreviousHistoryId(id domain.NullableInt32)
	SetFileVersion(version int32)
	*T
}

/**
 * Generic updater to track changes to a given data type.
 *
 * It is used to generically track changes in loaders so they can later apply it to the database.
 */

type Updater[K comparable, T UpdaterEntry[K], P UpdaterEntryPointer[K, T]] struct {
	CurrentValues     map[K]P
	DirtyValues       map[K]P
	ValuesToInsert    []T
	ValuesToUpdate    []T
	targetFileVersion int32
}

func NewUpdater[K comparable, T UpdaterEntry[K], P UpdaterEntryPointer[K, T]](currentValues []T, targetFileVersion int32) *Updater[K, T, P] {
	currentValuesHash := make(map[K]P)
	for _, m := range currentValues {
		currentValuesHash[m.GetId()] = &m
	}

	return &Updater[K, T, P]{
		DirtyValues:       make(map[K]P),
		CurrentValues:     currentValuesHash,
		targetFileVersion: targetFileVersion,
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

		// Append then build the pointer over the list item
		// otherwise, we will get 2 different copies of the element.
		u.ValuesToUpdate = append(u.ValuesToUpdate, newMap)
		newMapPtr := P(&u.ValuesToUpdate[len(u.ValuesToUpdate)-1])

		// This line is weird but makes sense. at this point newMap = *m, thus GetHistoryID
		// contains the original HistoryID, and we want it to become PreviousHistoryID
		newMapPtr.SetPreviousHistoryId(newMap.GetHistoryId())
		newMapPtr.SetHistoryId(domain.NewNullableNullInt32())
		newMapPtr.SetFileVersion(u.targetFileVersion)
		u.DirtyValues[key] = newMapPtr
		return newMapPtr
	}

	var newMapData T

	// Append then build the pointer over the list item
	// otherwise, we will get 2 different copies of the element.
	u.ValuesToInsert = append(u.ValuesToInsert, newMapData)
	newMap := P(&u.ValuesToInsert[len(u.ValuesToInsert)-1])

	newMap.SetId(key)
	newMap.SetHistoryId(domain.NewNullableNullInt32())
	newMap.SetFileVersion(u.targetFileVersion)
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
