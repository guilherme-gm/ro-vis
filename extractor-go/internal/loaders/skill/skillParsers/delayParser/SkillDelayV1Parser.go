package delayParser

import (
	"fmt"
	"regexp"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

type SkillDelayV1Parser struct {
}

func NewSkillDelayV1Parser() SkillDelayParser {
	return &SkillDelayV1Parser{}
}

func (p SkillDelayV1Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2010, time.April, 14, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC))
}

func (p SkillDelayV1Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		SkillDelayListV1Regex,
	}
}

func (p SkillDelayV1Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p SkillDelayV1Parser) ParseFile(filePath string, skillIdTable map[string]int) []rostructs.SkillDelayV1 {
	stringReencoder := decoders.ConvertEucKrToUtf8

	var skills []rostructs.SkillDelayV1
	decoder := decoders.NewLuaDecoder(stringReencoder)

	if err := validateFileConstats(filePath, SkFlags); err != nil {
		panic(err)
	}

	for flag, value := range SkFlags {
		decoder.CreateIntVar(flag, value)
	}

	decoder.CreateIntTable("SKID", skillIdTable)

	result := decoder.DecodeLuaTable(filePath, "SKILL_DELAY_LIST", &skills)
	if len(result.NotConsumedPaths) > 0 {
		fmt.Println("Not all keys were consumed.", result.NotConsumedPaths)
		panic("Not all keys were consumed.")
	}

	return skills
}

func (p SkillDelayV1Parser) Parse(basePath string, change *domain.UpdateChange, skillIdTable map[string]int) []rostructs.SkillDelayV1 {
	return p.ParseFile(basePath+"/"+change.Patch+"/"+change.File, skillIdTable)
}
