package repository

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type FromToRecord[T any] struct {
	LastUpdate domain.NullableString
	From       *domain.Record[T]
	To         *domain.Record[T]
}

type Repository struct {
	ItemRepository             *ItemRepository
	LoaderControllerRepository *LoaderControllerRepository
	PatchRepository            *PatchRepository
	QuestRepository            *QuestRepository
}

// NewRepository creates a new repository instance with the provided database connection
func NewRepository(db *database.Database) *Repository {
	return &Repository{
		ItemRepository:             NewItemRepository(db),
		LoaderControllerRepository: NewLoaderControllerRepository(db),
		PatchRepository:            NewPatchRepository(db),
		QuestRepository:            NewQuestRepository(db),
	}
}
