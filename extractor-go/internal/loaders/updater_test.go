package loaders_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
)

type testEntry struct {
	id                int
	historyId         domain.NullableInt32
	previousHistoryId domain.NullableInt32
	fileVersion       int32
	payloadA          string
	payloadB          string
}

// Satisfy UpdaterEntry
func (t testEntry) GetId() int                         { return t.id }
func (t testEntry) GetHistoryId() domain.NullableInt32 { return t.historyId }

// Satisfy UpdaterEntryPointer
func (t *testEntry) SetId(id int)                                 { t.id = id }
func (t *testEntry) SetHistoryId(id domain.NullableInt32)         { t.historyId = id }
func (t *testEntry) SetPreviousHistoryId(id domain.NullableInt32) { t.previousHistoryId = id }
func (t *testEntry) SetFileVersion(version int32)                 { t.fileVersion = version }

func TestNewUpdater_SeedsCurrentValues(t *testing.T) {
	seed := []testEntry{
		{id: 1, historyId: domain.NewNullableInt32(11), fileVersion: 1, payloadA: "a"},
		{id: 2, historyId: domain.NewNullableInt32(22), fileVersion: 1, payloadA: "b"},
	}
	u := loaders.NewUpdater(seed, 9)

	assert.Equal(t, 0, len(u.ValuesToInsert))
	assert.Equal(t, 0, len(u.ValuesToUpdate))
	assert.Equal(t, 0, len(u.DirtyValues))

	v1, ok1 := u.GetForRead(1)
	v2, ok2 := u.GetForRead(2)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.Equal(t, 1, v1.GetId())
	assert.Equal(t, 2, v2.GetId())
}

func TestGetForEdit_InsertNew(t *testing.T) {
	u := loaders.NewUpdater[int, testEntry](nil, 5)

	p := u.GetForEdit(10)
	assert.Equal(t, 10, p.GetId())
	assert.False(t, p.GetHistoryId().Valid)
	assert.Equal(t, int32(5), p.fileVersion)

	// Insert list should have one entry with id 10
	if assert.Equal(t, 1, len(u.ValuesToInsert)) {
		ins := u.ValuesToInsert[0]
		assert.Equal(t, 10, ins.GetId())
		assert.False(t, ins.GetHistoryId().Valid)
	}

	// Dirty should track editable pointer
	if got, ok := u.DirtyValues[10]; assert.True(t, ok) {
		assert.Equal(t, p, got)
	}
}

func TestGetForEdit_UpdateNew(t *testing.T) {
	u := loaders.NewUpdater[int, testEntry](nil, 5)

	p := u.GetForEdit(10)
	assert.Equal(t, 10, p.GetId())
	assert.False(t, p.GetHistoryId().Valid)
	assert.Equal(t, int32(5), p.fileVersion)
	p.payloadA = "newA"

	// Insert list should have one entry with id 10
	if assert.Equal(t, 1, len(u.ValuesToInsert)) {
		ins := u.ValuesToInsert[0]
		assert.Equal(t, 10, ins.GetId())
		assert.False(t, ins.GetHistoryId().Valid)
		assert.Equal(t, "newA", ins.payloadA)
	}

	// Changing it again should update the dirty value
	p1 := u.GetForEdit(10)
	p1.payloadA = "newB"
	if assert.Equal(t, 1, len(u.ValuesToInsert)) {
		ins := u.ValuesToInsert[0]
		assert.Equal(t, "newB", ins.payloadA)
	}

	// It should not list as value to update
	assert.Equal(t, 0, len(u.ValuesToUpdate))
}

func TestGetForEdit_UpdateExisting(t *testing.T) {
	seed := []testEntry{{id: 7, historyId: domain.NewNullableInt32(77), payloadA: "old"}}
	u := loaders.NewUpdater(seed, 6)

	p := u.GetForEdit(7)
	p.payloadA = "newA"
	// Editing existing should clone into update list with HistoryId cleared and PreviousHistoryId set
	if assert.Equal(t, 1, len(u.ValuesToUpdate)) {
		upd := u.ValuesToUpdate[0]
		assert.Equal(t, 7, upd.GetId())
		assert.False(t, upd.GetHistoryId().Valid)
		assert.Equal(t, "newA", upd.payloadA)
	}
	assert.Equal(t, domain.NewNullableInt32(77), p.previousHistoryId)
	assert.False(t, p.historyId.Valid)
	assert.Equal(t, int32(6), p.fileVersion)

	// Dirty should override reads
	read, ok := u.GetForRead(7)
	assert.True(t, ok)
	assert.False(t, read.GetHistoryId().Valid)
	assert.Equal(t, "newA", read.payloadA)

	// Edit again
	p1 := u.GetForEdit(7)
	p1.payloadB = "newB"
	if assert.Equal(t, 1, len(u.ValuesToUpdate)) {
		upd := u.ValuesToUpdate[0]
		assert.Equal(t, "newA", upd.payloadA)
		assert.Equal(t, "newB", upd.payloadB)
	}
}

func TestGetForRead_PrefersDirtyOverCurrent(t *testing.T) {
	seed := []testEntry{{id: 3, historyId: domain.NewNullableInt32(33)}}
	u := loaders.NewUpdater(seed, 1)

	p := u.GetForEdit(3)
	p.payloadA = "changed"

	read, ok := u.GetForRead(3)
	assert.True(t, ok)
	assert.Equal(t, "changed", read.payloadA)
}

func TestForEach_YieldsEachKeyOnce(t *testing.T) {
	seed := []testEntry{{id: 1}, {id: 2}}
	u := loaders.NewUpdater(seed, 0)
	_ = u.GetForEdit(2) // mark 2 as dirty
	_ = u.GetForEdit(3) // new

	seen := make(map[int]int)
	u.ForEach(func(k int, v testEntry) { seen[k]++ })

	assert.Equal(t, 1, seen[1])
	assert.Equal(t, 1, seen[2])
	assert.Equal(t, 1, seen[3])
	assert.Equal(t, 3, len(seen))
}
