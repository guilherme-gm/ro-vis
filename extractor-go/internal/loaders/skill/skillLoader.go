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
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/delayParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/descriptParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/idParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/infoParser"
)

type SkillLoader struct {
	repository   *repository.SkillRepository
	server       *server.Server
	skillParsers []skillParsers.SkillParser
}

// GetRelevantFiles returns a list of all files that are relevant to this loader's parsers.
func (l *SkillLoader) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		jobParsers.JobInheritListV3Regex,
		descriptParser.SkillDescriptV4Regex,
		idParser.SkillIdV2Regex,
		infoParser.SkillInfoListV2Regex,
		// infoParser.SkillInfoListV3Regex, // same as v2
		// infoParser.SkillInfoListV4Regex, // same as v2/v3
		// infoParser.SkillInfoListV5Regex, // same as v2/v3/v4
		delayParser.SkillDelayListV1Regex,
	}
}

func NewSkillLoader(server *server.Server) *SkillLoader {
	return &SkillLoader{
		repository: server.Repositories.SkillRepository,
		server:     server,
		skillParsers: []skillParsers.SkillParser{
			skillParsers.NewSkillV5Parser(),
			skillParsers.NewSkillV6Parser(),
			skillParsers.NewSkillV7Parser(),
			skillParsers.NewSkillV8Parser(),
		},
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
	if update.HasChangedAnyFiles(jobParsers.NewJobInehritListV3Parser().GetRelevantFiles()) {
		l.loadJobs(basePath, update, jobUpdater)
	}

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
	var targetParser skillParsers.SkillParser
	for _, skillParser := range l.skillParsers {
		if skillParser.IsUpdateInRange(&update) {
			targetParser = skillParser
		}
	}

	if targetParser == nil {
		panic("No parser found for update " + update.Name())
	}

	targetParser.Parse(basePath, update, skillUpdater, jobUpdater)

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
	change, err := update.GetChangeForFile(jobParsers.JobInehritListV3)
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

func (l *SkillLoader) Name() string {
	return "skill"
}
