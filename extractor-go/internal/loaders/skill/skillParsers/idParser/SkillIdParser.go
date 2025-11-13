package idParser

import (
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type SkillIdParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange) map[string]int
}
