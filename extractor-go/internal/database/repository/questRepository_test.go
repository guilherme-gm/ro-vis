package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockResult struct{ id int64 }

func (m mockResult) LastInsertId() (int64, error) { return m.id, nil }
func (m mockResult) RowsAffected() (int64, error) { return 1, nil }

func TestAddQuestsToHistory_Empty_NoOps(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewQuestRepository(mockDB)

	var quests []domain.Quest
	err := repo.AddQuestsToHistory(nil, "u1", &quests)

	require.NoError(t, err)
	mockDB.Dao.AssertExpectations(t)
}

func TestAddQuestsToHistory_Batching_And_Upsert(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewQuestRepository(mockDB)

	repo.BulkSize = 2

	// build 3 quests to force 2 batches (2 + 1)
	quests := make([]domain.Quest, 3)
	for i := range quests {
		quests[i] = domain.Quest{QuestID: int32(i + 1), FileVersion: 1}
	}

	// Expect two bulk inserts (slices of 2 and 1)
	mockDB.Dao.
		On("BulkInsertQuestHistory", mock.Anything, mock.MatchedBy(func(arg []dao.BulkInsertQuestHistoryParams) bool { return len(arg) == 2 })).
		Return(mockResult{id: 0}, nil).
		Once()
	mockDB.Dao.
		On("BulkInsertQuestHistory", mock.Anything, mock.MatchedBy(func(arg []dao.BulkInsertQuestHistoryParams) bool { return len(arg) == 1 })).
		Return(mockResult{id: 0}, nil).
		Once()

	// The repository will read ids for the update and upsert only those present in quests
	// It will later only Upsert the ones that are part of the current update
	ids := make([]dao.GetQuestsIdsInUpdateRow, 0, 3)
	ids = append(ids, dao.GetQuestsIdsInUpdateRow{HistoryID: 10, QuestID: 1})
	ids = append(ids, dao.GetQuestsIdsInUpdateRow{HistoryID: 20, QuestID: 2})
	ids = append(ids, dao.GetQuestsIdsInUpdateRow{HistoryID: 30, QuestID: 3})
	mockDB.Dao.
		On("GetQuestsIdsInUpdate", mock.Anything, "u1").
		Return(ids, nil).
		Twice()

	// We only upsert the ones that are part of the current bulk (1 and 2 / then 3)
	mockDB.Dao.
		On("BulkUpsertQuests", mock.Anything, mock.MatchedBy(func(arg []dao.BulkUpsertQuestParams) bool {
			if len(arg) != 2 {
				return false
			}
			return arg[0].QuestID == 1 && arg[0].HistoryID == 10 &&
				arg[1].QuestID == 2 && arg[1].HistoryID == 20
		})).
		Return(mockResult{id: 0}, nil).
		Once()

	mockDB.Dao.
		On("BulkUpsertQuests", mock.Anything, mock.MatchedBy(func(arg []dao.BulkUpsertQuestParams) bool {
			if len(arg) != 1 {
				return false
			}

			return arg[0].QuestID == 3 && arg[0].HistoryID == 30
		})).
		Return(mockResult{id: 0}, nil).
		Once()

	err := repo.AddQuestsToHistory(nil, "u1", &quests)
	require.NoError(t, err)
	mockDB.Dao.AssertExpectations(t)
}

func TestAddDeletedQuest_InsertsHistory_And_UpsertsDeleted(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewQuestRepository(mockDB)

	q := &domain.Quest{QuestID: 42, FileVersion: 2}
	mockDB.Dao.On("BulkInsertQuestHistory", mock.Anything, mock.MatchedBy(func(arg []dao.BulkInsertQuestHistoryParams) bool {
		return len(arg) == 1 && arg[0].QuestID == q.QuestID && arg[0].FileVersion == q.FileVersion
	})).Return(mockResult{id: 1234}, nil).Once()

	mockDB.Dao.On("UpsertQuest", mock.Anything, dao.UpsertQuestParams{
		QuestID:         42,
		LatestHistoryID: 1234,
		Deleted:         true,
	}).Return(mockResult{id: 0}, nil).Once()

	err := repo.AddDeletedQuest(nil, "upd", q)
	require.NoError(t, err)
	mockDB.Dao.AssertExpectations(t)
}

func TestCountChangesInUpdate(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewQuestRepository(mockDB)

	mockDB.Dao.On("CountChangedQuestsInUpdate", mock.Anything, "u").Return(int64(7), nil).Once()

	count, err := repo.CountChangesInUpdate(nil, "u")

	require.NoError(t, err)
	assert.Equal(t, 7, count)
	mockDB.Dao.AssertExpectations(t)
}

func TestGetChangesInUpdate_NoRows_ReturnsEmpty(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewQuestRepository(mockDB)

	mockDB.Dao.On("GetChangedQuests", mock.Anything, dao.GetChangedQuestsParams{Update: "u", Offset: 0, Limit: 10}).
		Return(nil, sql.ErrNoRows).Once()

	res, err := repo.GetChangesInUpdate(nil, "u", Pagination{Offset: 0, Limit: 10})

	require.NoError(t, err)
	assert.Len(t, res, 0)
	mockDB.Dao.AssertExpectations(t)
}

func TestGetQuestHistory_NoRows_ReturnsEmpty(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewQuestRepository(mockDB)

	mockDB.Dao.On("GetQuestHistory", mock.Anything, dao.GetQuestHistoryParams{QuestID: 99, Offset: 0, Limit: 5}).
		Return(nil, sql.ErrNoRows).Once()

	res, err := repo.GetQuestHistory(nil, 99, Pagination{Offset: 0, Limit: 5})
	require.NoError(t, err)
	assert.Len(t, res, 0)
	mockDB.Dao.AssertExpectations(t)
}

func TestGetQuests_ErrorIsIgnored_ReturnsEmpty(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewQuestRepository(mockDB)

	mockDB.Dao.On("GetQuestList", mock.Anything, dao.GetQuestListParams{Offset: 0, Limit: 3}).
		Return(nil, errors.New("db error")).Once()

	res, err := repo.GetQuests(nil, Pagination{Offset: 0, Limit: 3})
	require.NoError(t, err)
	assert.Len(t, res, 0)
	mockDB.Dao.AssertExpectations(t)
}

func TestCountQuests_CallsCountItems(t *testing.T) {
	mockDB := database.NewMockDatabase(t)
	repo := NewQuestRepository(mockDB)

	// Note: repository.CountQuests currently calls CountItems (not CountQuests)
	mockDB.Dao.On("CountItems", mock.Anything).Return(int64(123), nil).Once()

	count, err := repo.CountQuests(nil)
	require.NoError(t, err)
	assert.Equal(t, int32(123), count)
	mockDB.Dao.AssertExpectations(t)
}
