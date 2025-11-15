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
	// Don't trust this value as constant, I've found 1 case it is wrong and breaks everything
	// SillID always comes first anyway, we better trust it.
	// base.Constant = domain.NewNullableString(s.Constant)
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

	base.SkillScale = make([]domain.SkillScaleEntry, len(s.SkillScale))
	for i, v := range s.SkillScale {
		base.SkillScale[i] = domain.SkillScaleEntry{
			Level: int32(v.Level),
			X:     int32(v.X),
			Y:     int32(v.Y),
		}
	}

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

	base.SkillScale = make([]domain.SkillScaleEntry, len(s.SkillScale))
	for i, v := range s.SkillScale {
		base.SkillScale[i] = domain.SkillScaleEntry{
			Level: int32(v.Level),
			X:     int32(v.X),
			Y:     int32(v.Y),
		}
	}

	apCost := make([]int32, len(s.ApCost))
	for i, v := range s.ApCost {
		apCost[i] = int32(v)
	}
	base.ApCost = apCost

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
	base.SkillID = int32(s.SkillId)
	base.Name = domain.NewNullableString(s.SkillName)
	base.MaxLevel = domain.NewNullableInt32(int32(s.MaxLv))
	base.IsPassive = domain.NewNullableBool(s.IsPassive)

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

	base.SkillScale = make([]domain.SkillScaleEntry, len(s.SkillScale))
	for i, v := range s.SkillScale {
		base.SkillScale[i] = domain.SkillScaleEntry{
			Level: int32(v.Level),
			X:     int32(v.X),
			Y:     int32(v.Y),
		}
	}

	apCost := make([]int32, len(s.ApCost))
	for i, v := range s.ApCost {
		apCost[i] = int32(v)
	}
	base.ApCost = apCost

	return base
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
	base.CastFlags = make([]int32, len(s.SkillFlag))
	for i, v := range s.SkillFlag {
		base.CastFlags[i] = int32(v)
	}

	base.CastFixedDelay = make([]int32, len(s.SkillCastFixedDelay))
	for i, v := range s.SkillCastFixedDelay {
		base.CastFixedDelay[i] = int32(v)
	}

	base.CastStatDelay = make([]int32, len(s.SkillCastStatDelay))
	for i, v := range s.SkillCastStatDelay {
		base.CastStatDelay[i] = int32(v)
	}

	base.SinglePostDelay = make([]int32, len(s.SkillSinglePostDelay))
	for i, v := range s.SkillSinglePostDelay {
		base.SinglePostDelay[i] = int32(v)
	}

	base.GlobalPostDelay = make([]int32, len(s.SkillGlobalPostDelay))
	for i, v := range s.SkillGlobalPostDelay {
		base.GlobalPostDelay[i] = int32(v)
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
