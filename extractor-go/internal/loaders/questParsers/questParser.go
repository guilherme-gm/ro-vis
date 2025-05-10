package questParsers

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type QuestParser interface {
	IsPatchInRange(patch *domain.Patch) bool
	HasFiles(patch *domain.Patch) bool
	Parse(patchPath string) []domain.Quest
}
