package dao

import (
	"encoding/json"

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

func (q *GetCurrentSkillsRow) ToDomain() domain.Skill {
	var spAmount []int32
	if q.SpAmount != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.SpAmount, &spAmount)
	}

	var attackRange []int32
	if q.AttackRange != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.AttackRange, &attackRange)
	}

	var needSkillList []domain.NeedSkillEntry
	if q.NeedSkillList != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.NeedSkillList, &needSkillList)
	}

	return domain.Skill{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         ToNullableInt32(q.HistoryID),
		SkillID:           q.SkillID,
		FileVersion:       q.FileVersion,
		Constant:          domain.NullableString(q.Constant),
		Name:              domain.NullableString(q.Name),
		MaxLevel:          domain.NullableInt32(q.MaxLevel),
		SpAmount:          spAmount,
		SeparateLevel:     domain.NullableBool(q.SeparateLevel),
		AttackRange:       attackRange,
		NeedSkillList:     needSkillList,
	}
}
