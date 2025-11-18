package rostructs

import (
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/utils"
)

// V1 was used between 2009-10-07 and before 2010-04-14
// It used the following files:
// - data/lua files/skillinfo/jobinheritlist.lub
// - data/lua files/skillinfo/skillinfo_f.lub
// - data/lua files/skillinfo/skilltreeview.lub
//
// The details are unknown
// type SkillV1 struct {
// 	SkillId int
// }

// V2 started on 2010-04-14 up to present
// It uses the following files:
// - data/lua files/skillinfoz/jobinheritlist.lub
// - data/lua files/skillinfoz/skilldescript.lub
// - data/lua files/skillinfoz/skillid.lub
// - data/lua files/skillinfoz/skillinfolist.lub
// - data/lua files/skillinfoz/skillinfo_f.lub
// - data/lua files/skillinfoz/skilltreeview.lub
//
// Eventually moved to data/luafiles514/lua files/skillinfoz/*

type RequiredSkillV2 struct {
	SkillId int `lua:"$$numeric:1"`
	Lv      int `lua:"$$numeric:2"`
}

type JobRequiredSkillV2 struct {
	Job            int               `lua:"@index"`
	RequiredSkills []RequiredSkillV2 `lua:"@sliceValue"`
}

type SkillInfoV2 struct {
	SkillId           int    `lua:"@index"`
	Constant          string `lua:"$$numeric:1"`
	SkillName         string
	MaxLv             int
	Type              string
	SpCost            []int `lua:"SpAmount"`
	CanSelectLevel    bool  `lua:"bSeperateLv"`
	AttackRange       []int
	RequiredSkills    []RequiredSkillV2    `lua:"_NeedSkillList"`
	JobRequiredSkills []JobRequiredSkillV2 `lua:"NeedSkillList"`
}

func (s SkillInfoV2) ToSkill(base domain.Skill) domain.Skill {
	// Don't trust this value as constant, I've found 1 case it is wrong and breaks everything
	// SillID always comes first anyway, we better trust it.
	// base.Constant = domain.NewNullableString(s.Constant)
	base = setBasicFields(base, s.SkillId, s.SkillName, s.MaxLv)
	base = setSpAndRange(base, s.SpCost, s.CanSelectLevel, s.AttackRange)
	base = setRequiredSkills(base, s.RequiredSkills, s.JobRequiredSkills)
	return base
}

type SkillInfoScale struct {
	Level int `lua:"@index"`
	X     int `lua:"x"`
	Y     int `lua:"y"`
}

type SkillInfoV3 struct {
	SkillId           int    `lua:"@index"`
	Constant          string `lua:"$$numeric:1"`
	SkillName         string
	MaxLv             int
	Type              string
	SpCost            []int `lua:"SpAmount"`
	CanSelectLevel    bool  `lua:"bSeperateLv"`
	AttackRange       []int
	RequiredSkills    []RequiredSkillV2    `lua:"_NeedSkillList"`
	JobRequiredSkills []JobRequiredSkillV2 `lua:"NeedSkillList"`
	// New in V3
	SkillScale []SkillInfoScale `lua:"SkillScale"`
}

func (s SkillInfoV3) ToSkill(base domain.Skill) domain.Skill {
	// Don't trust this value as constant, I've found 1 case it is wrong and breaks everything
	// SillID always comes first anyway, we better trust it.
	// base.Constant = domain.NewNullableString(s.Constant)
	base = setBasicFields(base, s.SkillId, s.SkillName, s.MaxLv)
	base = setSpAndRange(base, s.SpCost, s.CanSelectLevel, s.AttackRange)
	base = setRequiredSkills(base, s.RequiredSkills, s.JobRequiredSkills)
	base.SkillScale = convertSkillScale(s.SkillScale)
	return base
}

type SkillInfoV4 struct {
	SkillId           int    `lua:"@index"`
	Constant          string `lua:"$$numeric:1"`
	SkillName         string
	MaxLv             int
	Type              string
	SpCost            []int `lua:"SpAmount"`
	CanSelectLevel    bool  `lua:"bSeperateLv"`
	AttackRange       []int
	RequiredSkills    []RequiredSkillV2    `lua:"_NeedSkillList"`
	JobRequiredSkills []JobRequiredSkillV2 `lua:"NeedSkillList"`
	SkillScale        []SkillInfoScale     `lua:"SkillScale"`
	// New in V4
	ApCost []int `lua:"ApAmount"`
}

func (s SkillInfoV4) ToSkill(base domain.Skill) domain.Skill {
	// Don't trust this value as constant, I've found 1 case it is wrong and breaks everything
	// SillID always comes first anyway, we better trust it.
	// base.Constant = domain.NewNullableString(s.Constant)
	base = setBasicFields(base, s.SkillId, s.SkillName, s.MaxLv)
	base = setSpAndRange(base, s.SpCost, s.CanSelectLevel, s.AttackRange)
	base = setRequiredSkills(base, s.RequiredSkills, s.JobRequiredSkills)
	base.SkillScale = convertSkillScale(s.SkillScale)
	base.ApCost = utils.ConvertIntSliceToInt32(s.ApCost)
	return base
}

type SkillInfoV5 struct {
	SkillId           int    `lua:"@index"`
	Constant          string `lua:"$$numeric:1"`
	SkillName         string
	MaxLv             int
	Type              string
	SpCost            []int `lua:"SpAmount"`
	CanSelectLevel    bool  `lua:"bSeperateLv"`
	AttackRange       []int
	RequiredSkills    []RequiredSkillV2    `lua:"_NeedSkillList"`
	JobRequiredSkills []JobRequiredSkillV2 `lua:"NeedSkillList"`
	SkillScale        []SkillInfoScale     `lua:"SkillScale"`
	ApCost            []int                `lua:"ApAmount"`
	// New in V5
	IsPassive bool
}

func (s SkillInfoV5) ToSkill(base domain.Skill) domain.Skill {
	// Don't trust this value as constant, I've found 1 case it is wrong and breaks everything
	// SillID always comes first anyway, we better trust it.
	// base.Constant = domain.NewNullableString(s.Constant)
	base = setBasicFields(base, s.SkillId, s.SkillName, s.MaxLv)
	base.IsPassive = domain.NewNullableBool(s.IsPassive)
	base = setSpAndRange(base, s.SpCost, s.CanSelectLevel, s.AttackRange)
	base = setRequiredSkills(base, s.RequiredSkills, s.JobRequiredSkills)
	base.SkillScale = convertSkillScale(s.SkillScale)
	base.ApCost = utils.ConvertIntSliceToInt32(s.ApCost)
	return base
}

func setBasicFields(base domain.Skill, skillId int, skillName string, maxLv int) domain.Skill {
	base.SkillID = int32(skillId)
	base.Name = domain.NewNullableString(skillName)
	base.MaxLevel = domain.NewNullableInt32(int32(maxLv))
	return base
}

func setSpAndRange(base domain.Skill, spCost []int, canSelect bool, attackRange []int) domain.Skill {
	base.SpCost = utils.ConvertIntSliceToInt32(spCost)
	base.CanSelectLevel = domain.NewNullableBool(canSelect)
	base.AttackRange = utils.ConvertIntSliceToInt32(attackRange)
	return base
}

func setRequiredSkills(base domain.Skill, req []RequiredSkillV2, jobReq []JobRequiredSkillV2) domain.Skill {
	base.RequiredSkills = make([]domain.NeedSkillEntry, len(req))
	for i, v := range req {
		base.RequiredSkills[i] = domain.NeedSkillEntry{
			SkillID: int32(v.SkillId),
			Level:   int32(v.Lv),
		}
	}

	base.JobRequiredSkills = make([]domain.JobRequiredSkillEntry, len(jobReq))
	for i, v := range jobReq {
		entry := domain.JobRequiredSkillEntry{
			JobId:  int32(v.Job),
			Skills: make([]domain.NeedSkillEntry, len(v.RequiredSkills)),
		}

		for j, w := range v.RequiredSkills {
			entry.Skills[j] = domain.NeedSkillEntry{
				SkillID: int32(w.SkillId),
				Level:   int32(w.Lv),
			}
		}

		base.JobRequiredSkills[i] = entry
	}

	return base
}

func convertSkillScale(src []SkillInfoScale) []domain.SkillScaleEntry {
	out := make([]domain.SkillScaleEntry, len(src))

	for i, v := range src {
		out[i] = domain.SkillScaleEntry{
			Level: int32(v.Level),
			X:     int32(v.X),
			Y:     int32(v.Y),
		}
	}

	return out
}

type SkillDelayV1 struct {
	SkillId              int `lua:"@index"`
	SkillFlag            []int
	SkillCastFixedDelay  []int
	SkillCastStatDelay   []int
	SkillSinglePostDelay []int
	SkillGlobalPostDelay []int
}

func (s SkillDelayV1) ToSkill(base domain.Skill) domain.Skill {
	base.CastFlags = utils.ConvertIntSliceToInt32(s.SkillFlag)
	base.CastFixedDelay = utils.ConvertIntSliceToInt32(s.SkillCastFixedDelay)
	base.CastStatDelay = utils.ConvertIntSliceToInt32(s.SkillCastStatDelay)
	base.SinglePostDelay = utils.ConvertIntSliceToInt32(s.SkillSinglePostDelay)
	base.GlobalPostDelay = utils.ConvertIntSliceToInt32(s.SkillGlobalPostDelay)
	return base
}

type SkillDescript struct {
	SkillId     int      `lua:"@index"`
	Description []string `lua:"@sliceValue"`
}

func (s SkillDescript) ToSkill(base domain.Skill) domain.Skill {
	base.Description = domain.NewNullableString(strings.Join(s.Description, "\n"))
	return base
}
