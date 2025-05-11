package questParsers

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type QuestParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	Parse(basePath string, update *domain.Update) []domain.Quest
}
