package infoParser

import (
	"fmt"
	"regexp"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

type SkillInfoV5Parser struct {
}

func NewSkillInfoV5Parser() SkillInfoV5Parser {
	return SkillInfoV5Parser{}
}

func (p SkillInfoV5Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2010, time.April, 14, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC))
}

func (p SkillInfoV5Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		SkillInfoListV5Regex,
	}
}

func (p SkillInfoV5Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p SkillInfoV5Parser) ParseFile(filePath string, jobIdTable map[string]int, skillIdTable map[string]int) []rostructs.SkillInfoV5 {
	stringReencoder := decoders.ConvertEucKrToUtf8

	var skills []rostructs.SkillInfoV5
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

func (p SkillInfoV5Parser) Parse(basePath string, change *domain.UpdateChange, jobIdTable map[string]int, skillIdTable map[string]int) []rostructs.SkillInfoV5 {
	return p.ParseFile(basePath+"/"+change.Patch+"/"+change.File, jobIdTable, skillIdTable)
}
