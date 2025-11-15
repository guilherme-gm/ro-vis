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
	var spCost []int32
	if q.SpCost != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.SpCost, &spCost)
	}

	var apCost []int32
	if q.ApCost != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.ApCost, &apCost)
	}

	var attackRange []int32
	if q.AttackRange != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.AttackRange, &attackRange)
	}

	var needSkillList []domain.NeedSkillEntry
	if q.RequiredSkills != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.RequiredSkills, &needSkillList)
	}

	var jobRequiredSkills []domain.JobRequiredSkillEntry
	if q.JobRequiredSkills != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.JobRequiredSkills, &jobRequiredSkills)
	}

	var skillScale []domain.SkillScaleEntry
	if q.SkillScale != nil {
		// Parse the JSON string into RewardItem slice
		json.Unmarshal(q.SkillScale, &skillScale)
	}

	return domain.Skill{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         ToNullableInt32(q.HistoryID),
		SkillID:           q.SkillID,
		FileVersion:       q.FileVersion,
		Constant:          domain.NullableString(q.Constant),
		Name:              domain.NullableString(q.Name),
		Description:       domain.NullableString(q.Description),
		MaxLevel:          domain.NullableInt32(q.MaxLevel),
		SpCost:            spCost,
		ApCost:            apCost,
		CanSelectLevel:    domain.NullableBool(q.CanSelectLevel),
		AttackRange:       attackRange,
		RequiredSkills:    needSkillList,
		JobRequiredSkills: jobRequiredSkills,
		SkillScale:        skillScale,
	}
}
