export type NeedSkillEntry = {
	SkillID: number;

	Level: number;
}

export type JobRequiredSkillEntry = {
	JobId: number;

	Skills: NeedSkillEntry[];
}

export type SkillScaleEntry = {
	Level: number;

	X: number;

	Y: number;
}

export type Skill = {
	HistoryID: string;

	FileVersion: number;

	SkillID: number;

	Constant: string;

	Name: string;

	Description: string;

	MaxLevel: number;

	IsPassive: boolean;

	SpCost: number[];

	ApCost: number[];

	CanSelectLevel: boolean;

	AttackRange: number[];

	RequiredSkills: NeedSkillEntry[];

	JobRequiredSkills: JobRequiredSkillEntry[];

	SkillScale: SkillScaleEntry[];

	CastFlags: number[];

	CastFixedDelay: number[];

	CastStatDelay: number[];

	SinglePostDelay: number[];

	GlobalPostDelay: number[];

	Deleted: boolean;
}

export type MinSkill = {
	SkillID: number;

	LastUpdate: string;

	Constant: string;

	Name: string;
};
