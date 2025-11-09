package dao

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func (q *SkillsJob) ToDomain() domain.SkillJob {
	return domain.SkillJob{
		Constant:      q.Constant,
		JobId:         q.JobID,
		InheritedJob:  domain.NullableInt32(q.InheritedJob),
		InheritedJob2: domain.NullableInt32(q.InheritedJob2),
		FirstUpdate:   q.FirstUpdate,
		LastUpdate:    q.LastUpdate,
		Deleted:       q.Deleted,
	}
}
