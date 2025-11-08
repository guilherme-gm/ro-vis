-- Trimmed down version of data/luafiles514/lua files/skillinfoz/skillinfolist.lub
-- contains only a small set of data to be used in tests

SKILL_INFO_LIST = {
	[SKID.SN_WINDWALK] = {
		"SN_WINDWALK",
		SkillName = "Wind Walk",
		MaxLv = 10,
		SpAmount = { 46, 52, 58, 64, 70, 76, 82, 88, 94, 100 },
		bSeperateLv = true,
		AttackRange = { 1, 1, 1, 1, 1, 1, 1, 1, 1, 1 },
		_NeedSkillList = {
			{ SKID.AC_CONCENTRATION, 9 }
		}
	},
	[SKID.AL_RUWACH] = {
		"AL_RUWACH",
		SkillName = "Ruach",
		MaxLv = 1,
		SpAmount = { 10 },
		bSeperateLv = false,
		AttackRange = { 10 }
	},
	[SKID.WS_MELTDOWN] = {
		"WS_MELTDOWN",
		SkillName = "Melt Down",
		MaxLv = 10,
		SpAmount = { 50, 50, 60, 60, 70, 70, 80, 80, 90, 90 },
		bSeperateLv = true,
		AttackRange = { 1, 1, 1, 1, 1, 1, 1, 1, 1, 1 },
		_NeedSkillList = {
			{ SKID.BS_SKINTEMPER, 3 },
			{ SKID.BS_HILTBINDING, 1 },
			{ SKID.BS_WEAPONRESEARCH, 5 },
			{ SKID.BS_OVERTHRUST, 3 }
		}
	},
	[SKID.CG_ARROWVULCAN] = {
		"CG_ARROWVULCAN",
		SkillName = "Arrow Vulcan",
		MaxLv = 10,
		SpAmount = { 12, 14, 16, 18, 20, 22, 24, 26, 28, 30 },
		bSeperateLv = true,
		AttackRange = { 9, 9, 9, 9, 9, 9, 9, 9, 9, 9 },
		NeedSkillList = {
			[JOBID.JT_BARD_H] = {
				{ SKID.AC_DOUBLE, 5 },
				{ SKID.AC_SHOWER, 5 },
				{ SKID.BA_MUSICALSTRIKE, 1 }
			},
			[JOBID.JT_DANCER_H] = {
				{ SKID.AC_DOUBLE, 5 },
				{ SKID.AC_SHOWER, 5 },
				{ SKID.DC_THROWARROW, 1 }
			}
		}
	},
}
