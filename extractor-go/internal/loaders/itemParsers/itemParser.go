package itemParsers

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type ItemParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []string
	Parse(basePath string, update *domain.Update, existingDB map[int32]*domain.Item) []domain.Item
}
