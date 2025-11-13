package skill

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/jobParsers"
)

type SkillLoader struct {
	repository *repository.SkillRepository
	server     *server.Server
}

// GetRelevantFiles returns a list of all files that are relevant to this loader's parsers.
func (l *SkillLoader) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		jobParsers.JobInheritListV2Regex,
		skillIdV2Regex,
		skillInfoListV2Regex,
	}
}

func NewSkillLoader(server *server.Server) *SkillLoader {
	return &SkillLoader{
		repository: server.Repositories.SkillRepository,
		server:     server,
	}
}

func (l *SkillLoader) LoadPatch(tx *sql.Tx, basePath string, update domain.Update) {
	if !update.HasChangedAnyFiles(l.GetRelevantFiles()) {
		fmt.Println("Skipped - No meaningful file")
		return
	}

	fmt.Println("> Decoding...")

	existingJobs, err := l.repository.GetSkillJobs(tx)
	if err != nil {
		panic(err)
	}

	jobUpdater := loaders.NewUpdater(existingJobs)
	l.loadJobs(basePath, update, jobUpdater)

	if len(jobUpdater.ValuesToInsert) > 0 {
		fmt.Println("> Inserting new skill jobs... ", len(jobUpdater.ValuesToInsert))
		if err := l.repository.AddSkillJobs(tx, jobUpdater.ValuesToInsert); err != nil {
			panic(err)
		}
	}
	if len(jobUpdater.ValuesToUpdate) > 0 {
		fmt.Println("> Updating skill jobs... ", len(jobUpdater.ValuesToUpdate))
		if err := l.repository.AddSkillJobs(tx, jobUpdater.ValuesToUpdate); err != nil {
			panic(err)
		}
	}

	existingSkills, err := l.repository.GetCurrentSkills(tx)
	if err != nil {
		panic(err)
	}

	skillUpdater := loaders.NewUpdater(existingSkills)
	l.loadSkillIDs(basePath, update, skillUpdater)
	l.loadSkillInfos(basePath, update, skillUpdater, jobUpdater)
	l.loadSkillDesc(basePath, update, skillUpdater)

	if len(skillUpdater.ValuesToInsert) > 0 {
		fmt.Println("> Inserting new skills... ", len(skillUpdater.ValuesToInsert))
		res, err := l.repository.AddSkillsToHistory(tx, update.Name(), skillUpdater.ValuesToInsert)
		if err != nil {
			panic(err)
		}

		fmt.Println("\tResult: ", res)
	}

	if len(skillUpdater.ValuesToUpdate) > 0 {
		fmt.Println("> Updating skills... ", len(skillUpdater.ValuesToUpdate))
		res, err := l.repository.AddSkillsToHistory(tx, update.Name(), skillUpdater.ValuesToUpdate)
		if err != nil {
			panic(err)
		}

		fmt.Println("\tResult: ", res)
	}
}

func (l *SkillLoader) loadJobs(basePath string, update domain.Update, jobUpdater *loaders.Updater[string, domain.SkillJob, *domain.SkillJob]) {
	change, err := update.GetChangeForFile(jobParsers.JobInehritListV2)
	if err != nil {
		panic(err)
	}

	jobParser := jobParsers.NewJobInehritListV3Parser()
	jobList := jobParser.Parse(basePath, &change)

	fileJobMap := make(map[string]bool)
	for _, fileJob := range jobList {
		fileJobMap[fileJob.Constant] = true
		existingJob, exists := jobUpdater.GetForRead(fileJob.Constant)
		shouldSave := false

		if !exists || existingJob.JobId != int32(fileJob.JobId) || existingJob.InheritedJob != fileJob.InheritedJob || existingJob.InheritedJob2 != fileJob.InheritedJob2 {
			shouldSave = true
		}

		if shouldSave {
			newJob := jobUpdater.GetForEdit(fileJob.Constant)
			newJob.Constant = fileJob.Constant
			newJob.JobId = int32(fileJob.JobId)
			newJob.InheritedJob = fileJob.InheritedJob
			newJob.InheritedJob2 = fileJob.InheritedJob2
			if existingJob.FirstUpdate == "" {
				newJob.FirstUpdate = update.Name()
			}
			newJob.LastUpdate = update.Name()
			newJob.Deleted = false
		}
	}

	for _, existingJob := range jobUpdater.CurrentValues {
		if _, ok := fileJobMap[existingJob.Constant]; !ok {
			deletedJob := jobUpdater.GetForEdit(existingJob.Constant)
			deletedJob.Deleted = true
			deletedJob.LastUpdate = update.Name()
		}
	}
}

func (l *SkillLoader) loadSkillIDs(basePath string, update domain.Update, skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill]) {
	change, err := update.GetChangeForFile(skillIdV2)
	if err != nil {
		panic(err)
	}

	skillParser := NewSkillIdV2Parser()
	skillList := skillParser.Parse(basePath, &change)

	fileSkillMap := make(map[int32]bool)
	for fileConst, fileId := range skillList {
		fileSkillMap[int32(fileId)] = true
		existingSkill, exists := skillUpdater.GetForRead(int32(fileId))
		shouldSave := false

		if !exists || existingSkill.SkillID != int32(fileId) || existingSkill.Constant.String != fileConst {
			shouldSave = true
		}

		if shouldSave {
			newSkill := skillUpdater.GetForEdit(int32(fileId))
			newSkill.Constant = domain.NewNullableString(fileConst)
			newSkill.SkillID = int32(fileId)
		}
	}

	for _, existingSkill := range skillUpdater.CurrentValues {
		if _, ok := fileSkillMap[existingSkill.SkillID]; !ok {
			deletedSkill := skillUpdater.GetForEdit(existingSkill.SkillID)
			deletedSkill.Deleted = true
		}
	}
}

type KV[T comparable] struct {
	Key   T
	Value int
}

func (l *SkillLoader) loadSkillInfos(
	basePath string,
	update domain.Update,
	skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill],
	jobUpdater *loaders.Updater[string, domain.SkillJob, *domain.SkillJob],
) {
	change, err := update.GetChangeForFile(skillInfoListV2)
	if err != nil {
		panic(err)
	}

	jobMap := make(map[string]int)
	jobUpdater.ForEach(func(key string, value domain.SkillJob) {
		jobMap[value.Constant] = int(value.JobId)
	})

	skillIdMap := make(map[string]int)
	skillUpdater.ForEach(func(key int32, value domain.Skill) {
		skillIdMap[value.Constant.String] = int(value.SkillID)
	})

	skillParser := NewSkillInfoV2Parser()
	skillList := skillParser.Parse(basePath, &change, jobMap, skillIdMap)

	fileSkillMap := make(map[int32]bool)
	for _, fileSkill := range skillList {
		fileSkillMap[int32(fileSkill.SkillId)] = true
		existingSkill, exists := skillUpdater.GetForRead(int32(fileSkill.SkillId))
		shouldSave := false

		fileSkillDomain := fileSkill.ToSkill(existingSkill)
		if !exists || !existingSkill.Equals(fileSkillDomain) {
			shouldSave = true
		}

		if shouldSave {
			newSkill := skillUpdater.GetForEdit(int32(fileSkill.SkillId))
			*newSkill = fileSkillDomain
		}
	}

	for _, existingSkill := range skillUpdater.CurrentValues {
		if _, ok := fileSkillMap[existingSkill.SkillID]; !ok {
			deletedSkill := skillUpdater.GetForEdit(existingSkill.SkillID)
			deletedSkill.Deleted = true
		}
	}
}

func (l *SkillLoader) loadSkillDesc(
	basePath string,
	update domain.Update,
	skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill],
) {
	change, err := update.GetChangeForFile(skillDescriptV2)
	if err != nil {
		panic(err)
	}

	skillIdMap := make(map[string]int)
	skillUpdater.ForEach(func(key int32, value domain.Skill) {
		skillIdMap[value.Constant.String] = int(value.SkillID)
	})

	skillParser := NewSkillDescriptV2Parser()
	skillList := skillParser.Parse(basePath, &change, skillIdMap)

	fileSkillMap := make(map[int32]bool)
	for _, fileSkill := range skillList {
		fileSkillMap[int32(fileSkill.SkillId)] = true
		existingSkill, exists := skillUpdater.GetForRead(int32(fileSkill.SkillId))
		shouldSave := false

		fileSkillDomain := fileSkill.ToSkill(existingSkill)
		if !exists || !existingSkill.Equals(fileSkillDomain) {
			shouldSave = true
		}

		if shouldSave {
			newSkill := skillUpdater.GetForEdit(int32(fileSkill.SkillId))
			*newSkill = fileSkillDomain
		}
	}

	for _, existingSkill := range skillUpdater.CurrentValues {
		if _, ok := fileSkillMap[existingSkill.SkillID]; !ok {
			if existingSkill.Description.Valid {
				deletedSkill := skillUpdater.GetForEdit(existingSkill.SkillID)
				deletedSkill.Description = domain.NewNullableNullString()
			}
		}
	}
}

func (l *SkillLoader) Name() string {
	return "skill"
}
