package loaders

import (
	"database/sql"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type Loader interface {
	GetRelevantFiles() []string
	LoadPatch(tx *sql.Tx, basePath string, update domain.Update)
	Name() string
}
