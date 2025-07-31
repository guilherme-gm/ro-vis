package itemParsers

import (
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type ItemParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, update *domain.Update, existingDB map[int32]*domain.Item) []domain.Item
}
