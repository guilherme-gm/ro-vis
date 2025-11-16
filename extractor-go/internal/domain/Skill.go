package domain

import "slices"

type SkillJob struct {
	Constant      string
	JobId         int32
	InheritedJob  NullableInt32
	InheritedJob2 NullableInt32
	FirstUpdate   string
	LastUpdate    string
	Deleted       bool
}

func (s SkillJob) GetId() string {
	return s.Constant
}

func (s *SkillJob) SetId(id string) {
	s.Constant = id
}

func (s SkillJob) GetHistoryId() NullableInt32 {
	return NewNullableNullInt32() // SkillJob doesn't have HistoryID
}

func (s *SkillJob) SetHistoryId(id NullableInt32) {
	// SkillJob doesn't have HistoryID
}

func (s *SkillJob) SetPreviousHistoryId(id NullableInt32) {
	// SkillJob doesn't have PreviousHistoryID
}

type NeedSkillEntry struct {
	SkillID int32
	Level   int32
}

type JobRequiredSkillEntry struct {
	JobId  int32
	Skills []NeedSkillEntry
}

type SkillScaleEntry struct {
	Level int32
	X     int32
	Y     int32
}

type Skill struct {
	PreviousHistoryID NullableInt32
	HistoryID         NullableInt32
	SkillID           int32
	FileVersion       int32
	Constant          NullableString
	Name              NullableString
	Description       NullableString
	MaxLevel          NullableInt32
	IsPassive         NullableBool
	SpCost            []int32
	ApCost            []int32
	CanSelectLevel    NullableBool
	AttackRange       []int32
	RequiredSkills    []NeedSkillEntry
	JobRequiredSkills []JobRequiredSkillEntry
	SkillScale        []SkillScaleEntry
	CastFlags         []int32
	CastFixedDelay    []int32
	CastStatDelay     []int32
	SinglePostDelay   []int32
	GlobalPostDelay   []int32
	Deleted           bool
}

func (s Skill) GetId() int32 {
	return s.SkillID
}

func (s *Skill) SetId(id int32) {
	s.SkillID = id
}

func (s Skill) GetHistoryId() NullableInt32 {
	return s.HistoryID
}

func (s *Skill) SetHistoryId(id NullableInt32) {
	s.HistoryID = id
}

func (s *Skill) SetPreviousHistoryId(id NullableInt32) {
	s.PreviousHistoryID = id
}

// Checks if JobRequiredSkills are equal
// We consider equal if we have the same jobs and same requirements
func (s *Skill) isJobRequiredSkillEqual(otherSkill Skill) bool {
	if len(s.JobRequiredSkills) != len(otherSkill.JobRequiredSkills) {
		return false
	}

	for _, v := range s.JobRequiredSkills {
		var otherRequired JobRequiredSkillEntry
		otherRequired.JobId = -1
		// find matching job in otherSkill
		for _, otherJob := range otherSkill.JobRequiredSkills {
			if otherJob.JobId == v.JobId {
				otherRequired = otherJob
				break
			}
		}

		if otherRequired.JobId == -1 {
			return false
		}

		if len(v.Skills) != len(otherRequired.Skills) {
			return false
		}

		for _, w := range v.Skills {
			var otherSkill NeedSkillEntry
			otherSkill.SkillID = -1
			// find matching skill in otherRequired
			for _, otherRequiredSkill := range otherRequired.Skills {
				if otherRequiredSkill.SkillID == w.SkillID {
					otherSkill = otherRequiredSkill
					break
				}
			}

			if otherSkill.SkillID == -1 {
				return false
			}

			if w.SkillID != otherSkill.SkillID {
				return false
			}

			if w.Level != otherSkill.Level {
				return false
			}
		}
	}

	return true
}

func (s *Skill) Equals(otherSkill Skill) bool {
	// FileVersion is not checked, if the file has changed but the skill is the same, we don't care.
	return (s.SkillID == otherSkill.SkillID &&
		s.Constant == otherSkill.Constant &&
		s.Name == otherSkill.Name &&
		s.Description == otherSkill.Description &&
		s.MaxLevel == otherSkill.MaxLevel &&
		s.IsPassive == otherSkill.IsPassive &&
		slices.Equal(s.SpCost, otherSkill.SpCost) &&
		slices.Equal(s.ApCost, otherSkill.ApCost) &&
		s.CanSelectLevel == otherSkill.CanSelectLevel &&
		slices.Equal(s.AttackRange, otherSkill.AttackRange) &&
		slices.Equal(s.RequiredSkills, otherSkill.RequiredSkills) &&
		s.isJobRequiredSkillEqual(otherSkill) &&
		slices.Equal(s.SkillScale, otherSkill.SkillScale)) &&
		slices.Equal(s.CastFlags, otherSkill.CastFlags) &&
		slices.Equal(s.CastFixedDelay, otherSkill.CastFixedDelay) &&
		slices.Equal(s.CastStatDelay, otherSkill.CastStatDelay) &&
		slices.Equal(s.SinglePostDelay, otherSkill.SinglePostDelay) &&
		slices.Equal(s.GlobalPostDelay, otherSkill.GlobalPostDelay)
}

type MinSkill struct {
	SkillID    int32
	LastUpdate string
	Constant   NullableString
	Name       NullableString
}
