package dao

import (
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type parsedSkillArrays struct {
	SpCost            []int32
	ApCost            []int32
	AttackRange       []int32
	RequiredSkills    []domain.NeedSkillEntry
	JobRequiredSkills []domain.JobRequiredSkillEntry
	SkillScale        []domain.SkillScaleEntry
	CastFlags         []int32
	CastFixedDelay    []int32
	CastStatDelay     []int32
	SinglePostDelay   []int32
	GlobalPostDelay   []int32
}

func parseSkillArrays(
	spCost []byte,
	apCost []byte,
	attackRange []byte,
	requiredSkills []byte,
	jobRequiredSkills []byte,
	skillScale []byte,
	castFlags []byte,
	castFixedDelay []byte,
	castStatDelay []byte,
	singlePostDelay []byte,
	globalPostDelay []byte,
) parsedSkillArrays {
	var out parsedSkillArrays
	if spCost != nil {
		json.Unmarshal(spCost, &out.SpCost)
	}
	if apCost != nil {
		json.Unmarshal(apCost, &out.ApCost)
	}
	if attackRange != nil {
		json.Unmarshal(attackRange, &out.AttackRange)
	}
	if requiredSkills != nil {
		json.Unmarshal(requiredSkills, &out.RequiredSkills)
	}
	if jobRequiredSkills != nil {
		json.Unmarshal(jobRequiredSkills, &out.JobRequiredSkills)
	}
	if skillScale != nil {
		json.Unmarshal(skillScale, &out.SkillScale)
	}
	if castFlags != nil {
		json.Unmarshal(castFlags, &out.CastFlags)
	}
	if castFixedDelay != nil {
		json.Unmarshal(castFixedDelay, &out.CastFixedDelay)
	}
	if castStatDelay != nil {
		json.Unmarshal(castStatDelay, &out.CastStatDelay)
	}
	if singlePostDelay != nil {
		json.Unmarshal(singlePostDelay, &out.SinglePostDelay)
	}
	if globalPostDelay != nil {
		json.Unmarshal(globalPostDelay, &out.GlobalPostDelay)
	}
	return out
}

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
	parsed := parseSkillArrays(
		q.SpCost,
		q.ApCost,
		q.AttackRange,
		q.RequiredSkills,
		q.JobRequiredSkills,
		q.SkillScale,
		q.CastFlags,
		q.CastFixedDelay,
		q.CastStatDelay,
		q.SinglePostDelay,
		q.GlobalPostDelay,
	)

	return domain.Skill{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         ToNullableInt32(q.HistoryID),
		SkillID:           q.SkillID,
		FileVersion:       q.FileVersion,
		Constant:          domain.NullableString(q.Constant),
		Name:              domain.NullableString(q.Name),
		Description:       domain.NullableString(q.Description),
		MaxLevel:          domain.NullableInt32(q.MaxLevel),
		IsPassive:         domain.NullableBool(q.IsPassive),
		SpCost:            parsed.SpCost,
		ApCost:            parsed.ApCost,
		CanSelectLevel:    domain.NullableBool(q.CanSelectLevel),
		AttackRange:       parsed.AttackRange,
		RequiredSkills:    parsed.RequiredSkills,
		JobRequiredSkills: parsed.JobRequiredSkills,
		SkillScale:        parsed.SkillScale,
		CastFlags:         parsed.CastFlags,
		CastFixedDelay:    parsed.CastFixedDelay,
		CastStatDelay:     parsed.CastStatDelay,
		SinglePostDelay:   parsed.SinglePostDelay,
		GlobalPostDelay:   parsed.GlobalPostDelay,
	}
}

func (q *PreviousSkillHistoryVw) ToDomain() domain.Skill {
	parsed := parseSkillArrays(
		q.SpCost,
		q.ApCost,
		q.AttackRange,
		q.RequiredSkills,
		q.JobRequiredSkills,
		q.SkillScale,
		q.CastFlags,
		q.CastFixedDelay,
		q.CastStatDelay,
		q.SinglePostDelay,
		q.GlobalPostDelay,
	)

	return domain.Skill{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         domain.NullableInt32(q.HistoryID),
		SkillID:           q.SkillID.Int32,
		FileVersion:       q.FileVersion.Int32,
		Constant:          domain.NullableString(q.Constant),
		Name:              domain.NullableString(q.Name),
		Description:       domain.NullableString(q.Description),
		MaxLevel:          domain.NullableInt32(q.MaxLevel),
		IsPassive:         domain.NullableBool(q.IsPassive),
		SpCost:            parsed.SpCost,
		ApCost:            parsed.ApCost,
		CanSelectLevel:    domain.NullableBool(q.CanSelectLevel),
		AttackRange:       parsed.AttackRange,
		RequiredSkills:    parsed.RequiredSkills,
		JobRequiredSkills: parsed.JobRequiredSkills,
		SkillScale:        parsed.SkillScale,
		CastFlags:         parsed.CastFlags,
		CastFixedDelay:    parsed.CastFixedDelay,
		CastStatDelay:     parsed.CastStatDelay,
		SinglePostDelay:   parsed.SinglePostDelay,
		GlobalPostDelay:   parsed.GlobalPostDelay,
	}
}

func (q *SkillsHistory) ToDomain() domain.Skill {
	parsed := parseSkillArrays(
		q.SpCost,
		q.ApCost,
		q.AttackRange,
		q.RequiredSkills,
		q.JobRequiredSkills,
		q.SkillScale,
		q.CastFlags,
		q.CastFixedDelay,
		q.CastStatDelay,
		q.SinglePostDelay,
		q.GlobalPostDelay,
	)

	return domain.Skill{
		PreviousHistoryID: domain.NullableInt32(q.PreviousHistoryID),
		HistoryID:         ToNullableInt32(q.HistoryID),
		SkillID:           q.SkillID,
		Constant:          domain.NullableString(q.Constant),
		Name:              domain.NullableString(q.Name),
		Description:       domain.NullableString(q.Description),
		MaxLevel:          domain.NullableInt32(q.MaxLevel),
		IsPassive:         domain.NullableBool(q.IsPassive),
		SpCost:            parsed.SpCost,
		ApCost:            parsed.ApCost,
		CanSelectLevel:    domain.NullableBool(q.CanSelectLevel),
		AttackRange:       parsed.AttackRange,
		RequiredSkills:    parsed.RequiredSkills,
		JobRequiredSkills: parsed.JobRequiredSkills,
		SkillScale:        parsed.SkillScale,
		CastFlags:         parsed.CastFlags,
		CastFixedDelay:    parsed.CastFixedDelay,
		CastStatDelay:     parsed.CastStatDelay,
		SinglePostDelay:   parsed.SinglePostDelay,
		GlobalPostDelay:   parsed.GlobalPostDelay,
	}
}
