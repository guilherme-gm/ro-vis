package i18n

import (
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

/**
 * First version of I18n data
 */
type I18nV1Parser struct {
}

func NewI18nV1Parser() I18nParser {
	return &I18nV1Parser{}
}

func (p I18nV1Parser) IsUpdateInRange(update *domain.Update) bool {
	return true
}

func (p I18nV1Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		scFiles,
	}
}

func (p I18nV1Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p I18nV1Parser) Parse(basePath string, change *domain.UpdateChange) []domain.I18n {
	result, err := decoders.DecodeLangCsv(basePath + "/" + change.Patch + "/" + change.File)
	if err != nil {
		panic(err)
	}

	i18ns := make([]domain.I18n, len(result))
	for idx, val := range result {
		i18ns[idx] = domain.I18n{
			I18nId:        val.Id,
			ContainerFile: change.File,
			EnText:        val.EnText,
			PtBrText:      val.PtBrText,
			Deleted:       false,
			FileVersion:   1,
		}
	}

	return i18ns
}
