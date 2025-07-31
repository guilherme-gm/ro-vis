package loaders

import (
	"database/sql"
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type Loader interface {
	GetRelevantFiles() []*regexp.Regexp
	LoadPatch(tx *sql.Tx, basePath string, update domain.Update)
	Name() string
}
