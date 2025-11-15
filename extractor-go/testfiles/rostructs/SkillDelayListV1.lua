-- Trimmed down version of data/luafiles514/lua files/skillinfoz/SkillDelayList.lub
-- contains only a small set of data to be used in tests

SKILL_DELAY_LIST = {
	[SKID.AC_DOUBLE] = {
		SkillGlobalPostDelay = { 100, 100, 100, 100, 100, 100, 100, 100, 100, 100 }
	},
	[SKID.AC_SHOWER] = {
		SkillGlobalPostDelay = { 100, 100, 100, 100, 100, 100, 100, 100, 100, 100 }
	},
	[SKID.BA_MUSICALSTRIKE] = {
		SkillCastFixedDelay = { 0, 0, 0, 0, 0 },
		SkillCastStatDelay = { 500, 500, 500, 500, 500 },
		SkillSinglePostDelay = { 0, 0, 0, 0, 0 },
		SkillGlobalPostDelay = { 300, 300, 300, 300, 300 }
	},
	[SKID.DC_THROWARROW] = {
		SkillCastFixedDelay = { 0, 0, 0, 0, 0 },
		SkillCastStatDelay = { 500, 500, 500, 500, 500 },
		SkillSinglePostDelay = { 0, 0, 0, 0, 0 },
		SkillGlobalPostDelay = { 300, 300, 300, 300, 300 }
	},
	[SKID.SN_WINDWALK] = {
		SkillCastFixedDelay = { 500, 600, 700, 800, 900, 1000, 1100, 1200, 1300, 1400 },
		SkillCastStatDelay = { 1500, 1800, 2100, 2400, 2700, 3000, 3300, 3600, 3900, 4200 },
		SkillGlobalPostDelay = { 2000, 2000, 2000, 2000, 2000, 2000, 2000, 2000, 2000, 2000 }
	},
	[SKID.WS_MELTDOWN] = {
		SkillCastFixedDelay = { 5, 5, 6, 6, 7, 7, 8, 8, 9, 10 },
		SkillCastStatDelay = { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }
	},
	[SKID.CG_ARROWVULCAN] = {
		SkillCastFixedDelay = { 500, 500, 500, 500, 500, 500, 500, 500, 500, 500 },
		SkillCastStatDelay = { 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500 },
		SkillGlobalPostDelay = { 600, 600, 600, 600, 600, 600, 600, 600, 600, 600 }, -- customized for tests
		SkillSinglePostDelay = { 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500 }
	},
	[SKID.GD_CHARGESHOUT_BEATING] = {
		SkillFlag = { SKFLAG_NOREDUCT, SKFLAG_DISABLE_FIXEDCASTING_REDUCTION },
		SkillCastFixedDelay = { 1000 },
		SkillSinglePostDelay = { 900000 }
	},
}
