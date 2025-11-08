package skill

import (
	"regexp"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/decoders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type SkillJobId struct {
	Constant      string
	JobId         int
	InheritedJob  domain.NullableInt32
	InheritedJob2 domain.NullableInt32
}

func newSkillJobId(constant string, jobId int) *SkillJobId {
	return &SkillJobId{
		Constant:      constant,
		JobId:         jobId,
		InheritedJob:  domain.NewNullableNullInt32(),
		InheritedJob2: domain.NewNullableNullInt32(),
	}
}

type JobInehritListV2Parser struct {
}

func NewJobInehritListV2Parser() *JobInehritListV2Parser {
	return &JobInehritListV2Parser{}
}

func (p JobInehritListV2Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2010, time.April, 14, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC))
}

func (p JobInehritListV2Parser) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		jobInehritListV2Regex,
	}
}

func (p JobInehritListV2Parser) HasFiles(update *domain.Update) bool {
	return update.HasChangedAnyFiles(p.GetRelevantFiles())
}

func (p JobInehritListV2Parser) parseFile(filePath string) []SkillJobId {
	stringReencoder := decoders.ConvertEucKrToUtf8

	var jobIds struct {
		Values []struct {
			Key   string `lua:"@index"`
			Value int    `lua:"@plainValue"`
		} `lua:"@plain"`
	}

	decoder := decoders.NewLuaDecoder(stringReencoder)
	decoder.DecodeLuaTable(filePath, "JOBID", &jobIds)

	var jobInehritList struct {
		Values []struct {
			Key   int `lua:"@index"`
			Value int `lua:"@plainValue"`
		} `lua:"@plain"`
	}

	decoder = decoders.NewLuaDecoder(stringReencoder)
	decoder.DecodeLuaTable(filePath, "JOB_INHERIT_LIST", &jobInehritList)

	var jobInehritList2 struct {
		Values []struct {
			Key   int `lua:"@index"`
			Value int `lua:"@plainValue"`
		} `lua:"@plain"`
	}
	decoder = decoders.NewLuaDecoder(stringReencoder)
	decoder.DecodeLuaTable(filePath, "JOB_INHERIT_LIST2", &jobInehritList2)

	var jobInehritListMap = make(map[int]*SkillJobId)
	for _, v := range jobIds.Values {
		jobInehritListMap[v.Value] = newSkillJobId(v.Key, v.Value)
	}

	for _, v := range jobInehritList.Values {
		jobInehritListMap[v.Key].InheritedJob = domain.NewNullableInt32(int32(v.Value))
	}

	for _, v := range jobInehritList2.Values {
		jobInehritListMap[v.Key].InheritedJob2 = domain.NewNullableInt32(int32(v.Value))
	}

	var result []SkillJobId
	for _, v := range jobInehritListMap {
		result = append(result, *v)
	}
	return result
}

func (p JobInehritListV2Parser) Parse(basePath string, change *domain.UpdateChange) []SkillJobId {
	return p.parseFile(basePath + "/" + change.Patch + "/" + change.File)
}
