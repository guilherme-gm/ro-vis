package questParsers

import (
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type QuestParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, update *domain.Update) []domain.Quest
}
