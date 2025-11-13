package infoParser

import (
	"regexp"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

type SkillInfoV1Parser struct {
}

func NewSkillInfoV1Parser() SkillInfoParser {
	return &SkillInfoV1Parser{}
}

func (p SkillInfoV1Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2009, time.October, 7, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(2010, time.April, 14, 0, 0, 0, 0, time.UTC))
}

func (p SkillInfoV1Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		SkillInfoV1Regex,
	}
}

func (p SkillInfoV1Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p SkillInfoV1Parser) Parse(basePath string, change *domain.UpdateChange, jobIdTable map[string]int, skillIdTable map[string]int) []rostructs.SkillInfoV2 {
	panic("SkillInfo V1 is not implemented")
}
