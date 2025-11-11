package skill

import (
	"fmt"
	"regexp"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

type SkillDescriptV2Parser struct {
}

func NewSkillDescriptV2Parser() SkillDescriptParser {
	return &SkillDescriptV2Parser{}
}

func (p SkillDescriptV2Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2010, time.April, 14, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC))
}

func (p SkillDescriptV2Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		skillDescriptV2Regex,
	}
}

func (p SkillDescriptV2Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p SkillDescriptV2Parser) parseFile(filePath string, skillIdTable map[string]int) []rostructs.SkillDescript {
	stringReencoder := decoders.ConvertEucKrToUtf8

	var skills []rostructs.SkillDescript
	decoder := decoders.NewLuaDecoder(stringReencoder)
	decoder.CreateIntTable("SKID", skillIdTable)

	result := decoder.DecodeLuaTable(filePath, "SKILL_DESCRIPT", &skills)
	if len(result.NotConsumedPaths) > 0 {
		fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
		panic("Not all keys were consumed.")
	}

	return skills
}

func (p SkillDescriptV2Parser) Parse(basePath string, change *domain.UpdateChange, skillIdTable map[string]int) []rostructs.SkillDescript {
	return p.parseFile(basePath+"/"+change.GetFullPath(), skillIdTable)
}
