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

type NeedSkillEntry struct {
	SkillID int32
	Level   int32
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
	SpAmount          []int32
	SeparateLevel     NullableBool
	AttackRange       []int32
	NeedSkillList     []NeedSkillEntry
	Deleted           bool
}

func (s Skill) GetId() int32 {
	return s.SkillID
}

func (s *Skill) SetId(id int32) {
	s.SkillID = id
}

func (s *Skill) Equals(otherSkill Skill) bool {
	// FileVersion is not checked, if the file has changed but the skill is the same, we don't care.
	if len(s.NeedSkillList) != len(otherSkill.NeedSkillList) {
		return false
	}

	return (s.SkillID == otherSkill.SkillID &&
		s.Constant == otherSkill.Constant &&
		s.Name == otherSkill.Name &&
		s.Description == otherSkill.Description &&
		s.MaxLevel == otherSkill.MaxLevel &&
		slices.Equal(s.SpAmount, otherSkill.SpAmount) &&
		s.SeparateLevel == otherSkill.SeparateLevel &&
		slices.Equal(s.AttackRange, otherSkill.AttackRange) &&
		slices.Equal(s.NeedSkillList, otherSkill.NeedSkillList))
}

type MinSkill struct {
	SkillID    int32
	LastUpdate string
	Name       NullableString
}
