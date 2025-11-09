package skill

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type skillJobUpdater struct {
	currentJobs  map[string]*domain.SkillJob
	dirtyJobs    map[string]*domain.SkillJob
	jobsToInsert []*domain.SkillJob
	jobsToUpdate []*domain.SkillJob
}

func newSkillJobUpdater(currentJobs []domain.SkillJob) *skillJobUpdater {
	currentJobHash := make(map[string]*domain.SkillJob)
	for _, m := range currentJobs {
		currentJobHash[m.Constant] = &m
	}

	return &skillJobUpdater{
		dirtyJobs:   make(map[string]*domain.SkillJob),
		currentJobs: currentJobHash,
	}
}

func (u *skillJobUpdater) getForRead(jobConst string) (domain.SkillJob, bool) {
	if m, ok := u.dirtyJobs[jobConst]; ok {
		return *m, true
	}

	if m, ok := u.currentJobs[jobConst]; ok {
		return *m, true
	}

	return domain.SkillJob{}, false
}

func (u *skillJobUpdater) getForEdit(jobConst string) *domain.SkillJob {
	if m, ok := u.dirtyJobs[jobConst]; ok {
		return m
	}

	if m, ok := u.currentJobs[jobConst]; ok {
		newMap := *m
		u.jobsToUpdate = append(u.jobsToUpdate, &newMap)
		u.dirtyJobs[jobConst] = &newMap
		return &newMap
	}

	newMap := domain.SkillJob{
		Constant: jobConst,
	}
	u.jobsToInsert = append(u.jobsToInsert, &newMap)
	u.dirtyJobs[jobConst] = &newMap
	return &newMap
}
