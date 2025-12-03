package dao

import (
	"database/sql"
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

func (q GetCurrentSkillsRow) ToDomain() domain.Skill {
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

func (s *BulkInsertSkillHistoryParams) FillFromDomain(skill *domain.Skill, update string) {
	spCostJson := domain.NewNullableNullString()
	if len(skill.SpCost) > 0 {
		jsonBytes, _ := json.Marshal(skill.SpCost)
		spCostJson = domain.NewNullableString(string(jsonBytes))
	}

	apCostJson := domain.NewNullableNullString()
	if len(skill.ApCost) > 0 {
		jsonBytes, _ := json.Marshal(skill.ApCost)
		apCostJson = domain.NewNullableString(string(jsonBytes))
	}

	attackRangeJson := domain.NewNullableNullString()
	if len(skill.AttackRange) > 0 {
		jsonBytes, _ := json.Marshal(skill.AttackRange)
		attackRangeJson = domain.NewNullableString(string(jsonBytes))
	}

	needSkillListJson := domain.NewNullableNullString()
	if len(skill.RequiredSkills) > 0 {
		jsonBytes, _ := json.Marshal(skill.RequiredSkills)
		needSkillListJson = domain.NewNullableString(string(jsonBytes))
	}

	jobRequiredSkillsJson := domain.NewNullableNullString()
	if len(skill.JobRequiredSkills) > 0 {
		jsonBytes, _ := json.Marshal(skill.JobRequiredSkills)
		jobRequiredSkillsJson = domain.NewNullableString(string(jsonBytes))
	}

	skillScaleJson := domain.NewNullableNullString()
	if len(skill.SkillScale) > 0 {
		jsonBytes, _ := json.Marshal(skill.SkillScale)
		skillScaleJson = domain.NewNullableString(string(jsonBytes))
	}

	castFlagsJson := domain.NewNullableNullString()
	if len(skill.CastFlags) > 0 {
		jsonBytes, _ := json.Marshal(skill.CastFlags)
		castFlagsJson = domain.NewNullableString(string(jsonBytes))
	}

	castFixedDelayJson := domain.NewNullableNullString()
	if len(skill.CastFixedDelay) > 0 {
		jsonBytes, _ := json.Marshal(skill.CastFixedDelay)
		castFixedDelayJson = domain.NewNullableString(string(jsonBytes))
	}

	castStatDelayJson := domain.NewNullableNullString()
	if len(skill.CastStatDelay) > 0 {
		jsonBytes, _ := json.Marshal(skill.CastStatDelay)
		castStatDelayJson = domain.NewNullableString(string(jsonBytes))
	}

	singlePostDelayJson := domain.NewNullableNullString()
	if len(skill.SinglePostDelay) > 0 {
		jsonBytes, _ := json.Marshal(skill.SinglePostDelay)
		singlePostDelayJson = domain.NewNullableString(string(jsonBytes))
	}

	globalPostDelayJson := domain.NewNullableNullString()
	if len(skill.GlobalPostDelay) > 0 {
		jsonBytes, _ := json.Marshal(skill.GlobalPostDelay)
		globalPostDelayJson = domain.NewNullableString(string(jsonBytes))
	}

	s.PreviousHistoryID = sql.NullInt32(skill.PreviousHistoryID)
	s.SkillId = skill.SkillID
	s.FileVersion = skill.FileVersion
	s.Update = update

	if !skill.Deleted {
		s.Constant = sql.NullString(skill.Constant)
		s.Name = sql.NullString(skill.Name)
		s.Description = sql.NullString(skill.Description)
		s.MaxLevel = sql.NullInt32(skill.MaxLevel)
		s.IsPassive = sql.NullBool(skill.IsPassive)
		s.SpCost = sql.NullString(spCostJson)
		s.ApCost = sql.NullString(apCostJson)
		s.CanSelectLevel = sql.NullBool(skill.CanSelectLevel)
		s.AttackRange = sql.NullString(attackRangeJson)
		s.RequiredSkills = sql.NullString(needSkillListJson)
		s.JobRequiredSkills = sql.NullString(jobRequiredSkillsJson)
		s.SkillScale = sql.NullString(skillScaleJson)
		s.CastFlags = sql.NullString(castFlagsJson)
		s.CastFixedDelay = sql.NullString(castFixedDelayJson)
		s.CastStatDelay = sql.NullString(castStatDelayJson)
		s.SinglePostDelay = sql.NullString(singlePostDelayJson)
		s.GlobalPostDelay = sql.NullString(globalPostDelayJson)
	}
}

func (s *BulkUpsertSkillParams) Fill(id int32, historyId int32, deleted bool) {
	s.SkillId = id
	s.HistoryID = historyId
	s.Deleted = deleted
}
