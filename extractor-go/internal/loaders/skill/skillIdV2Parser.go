package skill

import (
	"regexp"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type SkillIdV2Parser struct {
}

func NewSkillIdV2Parser() *SkillIdV2Parser {
	return &SkillIdV2Parser{}
}

func (p SkillIdV2Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2010, time.April, 14, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC))
}

func (p SkillIdV2Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		skillIdV2Regex,
	}
}

func (p SkillIdV2Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p SkillIdV2Parser) parseFile(filePath string) map[string]int {
	stringReencoder := decoders.ConvertEucKrToUtf8

	var skillIds struct {
		Values []struct {
			Key   string `lua:"@index"`
			Value int    `lua:"@plainValue"`
		} `lua:"@plain"`
	}

	decoder := decoders.NewLuaDecoder(stringReencoder)
	decoder.DecodeLuaTable(filePath, "SKID", &skillIds)

	var skillIdMap = make(map[string]int)
	for _, v := range skillIds.Values {
		skillIdMap[v.Key] = v.Value
	}

	return skillIdMap
}

func (p SkillIdV2Parser) Parse(basePath string, change *domain.UpdateChange) map[string]int {
	return p.parseFile(basePath + "/" + change.Patch + "/" + change.File)
}
