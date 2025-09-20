package repository

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildDAOPatch(id int32, name string, date time.Time, files []string) dao.Patch {
	b, _ := json.Marshal(files)
	return dao.Patch{ID: id, Name: name, Date: date, Files: b, Status: dao.PatchesStatusPending}
}

func TestListUpdates_PaginateAll_GroupingAndOverride(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	dateA := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	dateB := time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC)

	mockDB.Dao.On("ListPatches", mock.Anything, dateB).Return([]dao.Patch{
		buildDAOPatch(1, "patch1", dateA, []string{"a.txt", "b.txt"}),
		buildDAOPatch(2, "patch2", dateA, []string{"a.txt", "c.txt"}),
		buildDAOPatch(3, "patch3", dateB, []string{"a.txt", "d.txt"}),
	}, nil).Once()

	repo := NewPatchRepository(mockDB)
	updates, err := repo.ListUpdates(nil, dateB, PaginateAll)
	require.NoError(t, err)

	require.Len(t, updates, 2)
	assert.Equal(t, dateA, updates[0].Date, "The first update should have date A (oldest)")
	assert.Equal(t, dateB, updates[1].Date, "The second update should have date B (newest)")

	expected := []domain.UpdateChange{
		{Patch: "patch2", File: "a.txt"},
		{Patch: "patch1", File: "b.txt"},
		{Patch: "patch2", File: "c.txt"},
	}
	assert.ElementsMatch(t, expected, updates[0].Changes, "When 2 patches overlaps, it should use the newest one.")

	expected = []domain.UpdateChange{
		{Patch: "patch3", File: "a.txt"},
		{Patch: "patch3", File: "d.txt"},
	}
	assert.ElementsMatch(t, expected, updates[1].Changes, "The second update should only contain files from patch 3")
}

func TestListUpdates_NoPatches(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mockDB.Dao.On("ListPatches", mock.Anything, date).Return([]dao.Patch{}, nil).Once()
	repo := NewPatchRepository(mockDB)
	updates, err := repo.ListUpdates(nil, date, PaginateAll)
	require.NoError(t, err)
	assert.Len(t, updates, 0)
	mockDB.Dao.AssertExpectations(t)
}
