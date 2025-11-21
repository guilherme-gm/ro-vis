package repository

import (
	"database/sql"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type FromToRecord[T any] struct {
	LastUpdate domain.NullableString
	From       *domain.Record[T]
	To         *domain.Record[T]
}

// BaseRepository provides shared DB plumbing for repositories
// and a dao helper bound to an optional transaction.
type BaseRepository struct {
	DB database.IDatabase
}

func (r BaseRepository) dao(tx *sql.Tx) dao.IDAO {
	return r.DB.GetDAO(tx)
}

// chunkIndices returns [start,end) index pairs for splitting a list into batches of size chunkSize.
func chunkIndices(total int, chunkSize int) [][2]int {
	if chunkSize <= 0 || total <= 0 {
		return nil
	}
	segments := make([][2]int, 0, (total+chunkSize-1)/chunkSize)
	for start := 0; start < total; start += chunkSize {
		end := start + chunkSize
		if end > total {
			end = total
		}
		segments = append(segments, [2]int{start, end})
	}
	return segments
}

// insertDeletionAndUpsert orchestrates the common delete flow used by repositories:
// 1) insert a history row, 2) capture LastInsertId and set the entity history ID,
// 3) upsert the current table with Deleted=true.
func (r BaseRepository) insertDeletionAndUpsert(
	tx *sql.Tx,
	insertHistory func(dao.IDAO) (sql.Result, error),
	setHistoryID func(int64),
	upsert func(dao.IDAO) (sql.Result, error),
) error {
	q := r.dao(tx)
	res, err := insertHistory(q)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	setHistoryID(id)
	_, err = upsert(q)
	return err
}

type Repository struct {
	ItemRepository             *ItemRepository
	LoaderControllerRepository *LoaderControllerRepository
	PatchRepository            *PatchRepository
	QuestRepository            *QuestRepository
	I18nRepository             *I18nRepository
	MapRepository              *MapRepository
	SkillRepository            *SkillRepository
}

// NewRepository creates a new repository instance with the provided database connection
func NewRepository(db database.IDatabase) *Repository {
	return &Repository{
		ItemRepository:             NewItemRepository(db),
		LoaderControllerRepository: NewLoaderControllerRepository(db),
		PatchRepository:            NewPatchRepository(db),
		QuestRepository:            NewQuestRepository(db),
		I18nRepository:             NewI18nRepository(db),
		MapRepository:              NewMapRepository(db),
		SkillRepository:            NewSkillRepository(db),
	}
}
