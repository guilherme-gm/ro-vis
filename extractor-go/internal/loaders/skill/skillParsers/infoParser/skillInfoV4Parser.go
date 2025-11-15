package infoParser

import (
	"fmt"
	"regexp"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

type SkillInfoV4Parser struct {
}

func NewSkillInfoV4Parser() SkillInfoV4Parser {
	return SkillInfoV4Parser{}
}

func (p SkillInfoV4Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2010, time.April, 14, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC))
}

func (p SkillInfoV4Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		SkillInfoListV4Regex,
	}
}

func (p SkillInfoV4Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p SkillInfoV4Parser) ParseFile(filePath string, jobIdTable map[string]int, skillIdTable map[string]int) []rostructs.SkillInfoV4 {
	stringReencoder := decoders.ConvertEucKrToUtf8

	var skills []rostructs.SkillInfoV4
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

func (p SkillInfoV4Parser) Parse(basePath string, change *domain.UpdateChange, jobIdTable map[string]int, skillIdTable map[string]int) []rostructs.SkillInfoV4 {
	return p.ParseFile(basePath+"/"+change.Patch+"/"+change.File, jobIdTable, skillIdTable)
}
