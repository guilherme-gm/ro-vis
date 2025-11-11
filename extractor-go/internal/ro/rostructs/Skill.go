package rostructs

import (
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
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
	base.Constant = domain.NewNullableString(s.Constant)
	base.SkillID = int32(s.SkillId)
	base.Name = domain.NewNullableString(s.SkillName)
	base.MaxLevel = domain.NewNullableInt32(int32(s.MaxLv))

	spCost := make([]int32, len(s.SpCost))
	for i, v := range s.SpCost {
		spCost[i] = int32(v)
	}
	base.SpCost = spCost
	base.CanSelectLevel = domain.NewNullableBool(s.CanSelectLevel)
	base.AttackRange = make([]int32, len(s.AttackRange))
	for i, v := range s.AttackRange {
		base.AttackRange[i] = int32(v)
	}

	base.RequiredSkills = make([]domain.NeedSkillEntry, len(s.RequiredSkills))
	for i, v := range s.RequiredSkills {
		base.RequiredSkills[i] = domain.NeedSkillEntry{
			SkillID: int32(v.SkillId),
			Level:   int32(v.Lv),
		}
	}

	base.JobRequiredSkills = make([]domain.JobRequiredSkillEntry, len(s.JobRequiredSkills))
	for i, v := range s.JobRequiredSkills {
		base.JobRequiredSkills[i] = domain.JobRequiredSkillEntry{
			JobId:  int32(v.Job),
			Skills: make([]domain.NeedSkillEntry, len(v.RequiredSkills)),
		}

		for j, w := range v.RequiredSkills {
			base.JobRequiredSkills[i].Skills[j] = domain.NeedSkillEntry{
				SkillID: int32(w.SkillId),
				Level:   int32(w.Lv),
			}
		}
	}

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
