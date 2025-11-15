package skillParsers

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
)

type SkillParser interface {
	IsUpdateInRange(update *domain.Update) bool
	Parse(basePath string, update domain.Update, skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill], jobUpdater *loaders.Updater[string, domain.SkillJob, *domain.SkillJob])
}

type KV[T comparable] struct {
	Key   T
	Value int
}

func shouldSkipConst(constant string) bool {
	// These were present in a few 2015 updates of SkillID,
	// but they were not used and they conflict with other real skills
	// They were later removed in 2015-12.
	// They cause a lot of issues since we base on the numeric ID as the lead, and
	// no reason to fret over those 4, so they are skipped.
	if constant == "SR_FLASHCOMBO_ATK_STEP1" ||
		constant == "SR_FLASHCOMBO_ATK_STEP2" ||
		constant == "SR_FLASHCOMBO_ATK_STEP3" ||
		constant == "SR_FLASHCOMBO_ATK_STEP4" {
		return true
	}

	// Firstly introduced in 2015-12-23, which was fine, but later on it started conflicting with real skills
	// - NPC_LAST - on 2018-08-29 it started conflicting with new NPC skills
	//    - Just a marker, no reason to track
	// - RL_GLITTERING_GREED_ATK (2574) - on 2018 it started conflicting with SJ_LIGHTOFMOON (2574)
	//    - seems used for some internal handling (the skill does not exists)
	// - ELEMENTAL_LAST (8443) - on 2020-09-15+ conflicts with EM_EL_FLAMETECHNIC
	//    - Just a marker, no reason to track
	// so we are ignoring it.
	if constant == "NPC_LAST" ||
		constant == "RL_GLITTERING_GREED_ATK" ||
		constant == "ELEMENTAL_LAST" {
		return true
	}

	// These were added on 2016-05-04 and conflicts with real skills, but seems to be useless
	// As of 2025, they are no longer in kRO files
	// - WL_FREEZE_SP (2232)
	// - SC_STARTMARK (2284)
	// - SCRIPT_000 (11000)
	// It is unused
	if constant == "WL_ENDMARK" ||
		constant == "SC_STARTMARK" ||
		constant == "SCRIPT_000" {
		return true
	}

	return false
}
