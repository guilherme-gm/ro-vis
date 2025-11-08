package skill

import (
	"fmt"
	"regexp"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

type SkillInfoV2Parser struct {
}

func NewSkillInfoV2Parser() SkillInfoParser {
	return &SkillInfoV2Parser{}
}

func (p SkillInfoV2Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2010, time.April, 14, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC))
}

func (p SkillInfoV2Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		skillInfoListV2Regex,
	}
}

func (p SkillInfoV2Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p SkillInfoV2Parser) parseFile(filePath string, jobIdTable map[string]int, skillIdTable map[string]int) []rostructs.SkillInfoV2 {
	stringReencoder := decoders.ConvertEucKrToUtf8

	var skills []rostructs.SkillInfoV2
	decoder := decoders.NewLuaDecoder(stringReencoder)
	decoder.CreateIntTable("JOBID", jobIdTable)
	decoder.CreateIntTable("SKID", skillIdTable)

	result := decoder.DecodeLuaTable(filePath, "SKILL_INFO_LIST", &skills)
	if len(result.NotConsumedPaths) > 0 {
		fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
		panic("Not all keys were consumed.")
	}

	return skills
}

func (p SkillInfoV2Parser) Parse(basePath string, change *domain.UpdateChange, jobIdTable map[string]int, skillIdTable map[string]int) []rostructs.SkillInfoV2 {
	return p.parseFile(basePath+"/"+change.Patch+"/"+change.File, jobIdTable, skillIdTable)
}
