package repository

import (
	"database/sql"
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddSkillsToHistory_Empty_ReturnsZeroAndNoError(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewSkillRepository(mockDB)

	var skills []*domain.Skill
	res, err := repo.AddSkillsToHistory(nil, "u1", skills)

	require.NoError(t, err)
	assert.Equal(t, AddToHistoryResult{}, res)
	mockDB.Dao.AssertExpectations(t)
}

func TestAddSkillsToHistory_Batching_And_Counts(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewSkillRepository(mockDB)
	repo.BulkSize = 2

	// Build 3 skills to force two batches (2 + 1)
	skills := make([]*domain.Skill, 3)
	deletedCount := 0
	for i := range 3 {
		// alternate deleted flag
		isDel := i%3 == 0
		if isDel {
			deletedCount++
		}
		skills[i] = &domain.Skill{SkillID: int32(i + 1), FileVersion: 1, Deleted: isDel}
	}

	// Expect two bulk inserts (slices of 2 and 1)
	mockDB.Dao.
		On("BulkInsertSkillHistory", mock.Anything, mock.MatchedBy(func(arg []dao.BulkInsertSkillHistoryParams) bool { return len(arg) == 2 })).
		Return(mockResult{id: 0}, nil).
		Once()
	mockDB.Dao.
		On("BulkInsertSkillHistory", mock.Anything, mock.MatchedBy(func(arg []dao.BulkInsertSkillHistoryParams) bool { return len(arg) == 1 })).
		Return(mockResult{id: 0}, nil).
		Once()

	// IDs in update: generate one per skill
	ids := make([]dao.GetSkillsIdsInUpdateRow, 0, 3)
	for i := range 3 {
		ids = append(ids, dao.GetSkillsIdsInUpdateRow{HistoryID: int32(1000 + i + 1), SkillID: int32(i + 1)})
	}
	// Will be called once per batch
	mockDB.Dao.On("GetSkillsIdsInUpdate", mock.Anything, "u2").Return(ids, nil).Twice()

	// Expect upserts per batch
	mockDB.Dao.
		On("BulkUpsertSkills", mock.Anything, mock.MatchedBy(func(arg []dao.BulkUpsertSkillParams) bool { return len(arg) == 2 })).
		Return(mockResult{id: 0}, nil).
		Once()
	mockDB.Dao.
		On("BulkUpsertSkills", mock.Anything, mock.MatchedBy(func(arg []dao.BulkUpsertSkillParams) bool { return len(arg) == 1 })).
		Return(mockResult{id: 0}, nil).
		Once()

	res, err := repo.AddSkillsToHistory(nil, "u2", skills)
	require.NoError(t, err)

	// Verify counts
	assert.Equal(t, deletedCount, res.DeletedCount)
	assert.Equal(t, len(skills)-deletedCount, res.UpsertCount)

	mockDB.Dao.AssertExpectations(t)
}

func TestCountSkills_ReturnsValue(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewSkillRepository(mockDB)

	mockDB.Dao.On("CountSkills", mock.Anything).Return(int64(77), nil).Once()

	count, err := repo.CountSkills(nil)
	require.NoError(t, err)
	assert.Equal(t, int32(77), count)
	mockDB.Dao.AssertExpectations(t)
}

func TestGetSkills_ErrorIgnored_ReturnsEmpty(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewSkillRepository(mockDB)

	mockDB.Dao.On("GetSkillList", mock.Anything, dao.GetSkillListParams{Offset: 0, Limit: 3}).Return(nil, assert.AnError).Once()

	res, err := repo.GetSkills(nil, Pagination{Offset: 0, Limit: 3})
	require.NoError(t, err)
	assert.Len(t, res, 0)
	mockDB.Dao.AssertExpectations(t)
}

func TestGetSkillHistory_NoRows_ReturnsEmpty(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewSkillRepository(mockDB)

	mockDB.Dao.On("GetSkillHistory", mock.Anything, dao.GetSkillHistoryParams{SkillID: 10, Offset: 0, Limit: 5}).Return(nil, sql.ErrNoRows).Once()

	res, err := repo.GetSkillHistory(nil, 10, Pagination{Offset: 0, Limit: 5})
	require.NoError(t, err)
	assert.Len(t, res, 0)
	mockDB.Dao.AssertExpectations(t)
}

func TestSkill_GetChangesInUpdate_NoRows_ReturnsEmpty(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewSkillRepository(mockDB)

	mockDB.Dao.On("GetChangedSkills", mock.Anything, dao.GetChangedSkillsParams{Update: "upd", Offset: 0, Limit: 10}).Return(nil, sql.ErrNoRows).Once()

	res, err := repo.GetChangesInUpdate(nil, "upd", Pagination{Offset: 0, Limit: 10})
	require.NoError(t, err)
	assert.Len(t, res, 0)
	mockDB.Dao.AssertExpectations(t)
}

func TestCountChangesInUpdate_CallsDAO(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewSkillRepository(mockDB)

	mockDB.Dao.On("CountChangedSkillsInUpdate", mock.Anything, "u").Return(int64(5), nil).Once()

	count, err := repo.CountChangesInUpdate(nil, "u")
	require.NoError(t, err)
	assert.Equal(t, 5, count)
	mockDB.Dao.AssertExpectations(t)
}

func TestGetSkillJobs_NoRows_ReturnsEmpty(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewSkillRepository(mockDB)

	mockDB.Dao.On("GetSkillJobs", mock.Anything).Return(nil, sql.ErrNoRows).Once()

	res, err := repo.GetSkillJobs(nil)

	require.NoError(t, err)
	assert.Len(t, res, 0)
	mockDB.Dao.AssertExpectations(t)
}
