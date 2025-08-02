package i18n

import (
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type I18nParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange) []domain.I18n
}

var scFiles = regexp.MustCompile(`(?i)^data/i18n/sc/.*\.csv$`)
var scJsonFile = regexp.MustCompile(`(?i)^data/i18n/sc/sc.json$`)
