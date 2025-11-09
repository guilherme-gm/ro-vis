package skill

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

type SkillLoader struct {
	repository *repository.SkillRepository
	server     *server.Server
}

// GetRelevantFiles returns a list of all files that are relevant to this loader's parsers.
func (l *SkillLoader) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		jobInehritListV2Regex,
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

	jobUpdater := newSkillJobUpdater(*existingJobs)
	l.loadJobs(basePath, update, jobUpdater)

	if len(jobUpdater.jobsToInsert) > 0 {
		fmt.Println("> Inserting new skill jobs... ", len(jobUpdater.jobsToInsert))
		if err := l.repository.AddSkillJobs(tx, jobUpdater.jobsToInsert); err != nil {
			panic(err)
		}
	}
	if len(jobUpdater.jobsToUpdate) > 0 {
		fmt.Println("> Updating skill jobs... ", len(jobUpdater.jobsToUpdate))
		if err := l.repository.AddSkillJobs(tx, jobUpdater.jobsToUpdate); err != nil {
			panic(err)
		}
	}
}

func (l *SkillLoader) loadJobs(basePath string, update domain.Update, jobUpdater *skillJobUpdater) {
	change, err := update.GetChangeForFile(jobInehritListV2)
	if err != nil {
		panic(err)
	}

	jobParser := NewJobInehritListV2Parser()
	jobList := jobParser.Parse(basePath, &change)

	fileJobMap := make(map[string]bool)
	for _, fileJob := range jobList {
		fileJobMap[fileJob.Constant] = true
		existingJob, exists := jobUpdater.getForRead(fileJob.Constant)
		shouldSave := false

		if !exists || existingJob.JobId != int32(fileJob.JobId) || existingJob.InheritedJob != fileJob.InheritedJob || existingJob.InheritedJob2 != fileJob.InheritedJob2 {
			shouldSave = true
		}

		if shouldSave {
			newJob := jobUpdater.getForEdit(fileJob.Constant)
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

	for _, existingJob := range jobUpdater.currentJobs {
		if _, ok := fileJobMap[existingJob.Constant]; !ok {
			deletedJob := jobUpdater.getForEdit(existingJob.Constant)
			deletedJob.Deleted = true
			deletedJob.LastUpdate = update.Name()
		}
	}
}

func (l *SkillLoader) Name() string {
	return "skill"
}
