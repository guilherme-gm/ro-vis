package rostructs

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
