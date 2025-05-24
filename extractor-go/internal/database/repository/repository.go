package repository

import (
	"database/sql"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type FromToRecord[T any] struct {
	LastUpdate sql.NullString
	From       *domain.Record[T]
	To         *domain.Record[T]
}

type repositories struct {
	ItemRepository             *ItemRepository
	LoaderControllerRepository *LoaderControllerRepository
	patchRepository            *PatchRepository
	questRepository            *QuestRepository
}

var repositoriesCache repositories = repositories{}
