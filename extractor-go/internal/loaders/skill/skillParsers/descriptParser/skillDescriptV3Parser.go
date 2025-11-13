package descriptParser

import (
	"regexp"
)

var SkillDescriptV3 = "data/lua files/skillinfoz/skilldescript.lub"
var SkillDescriptV3Regex = regexp.MustCompile(`(?i)^` + SkillDescriptV3 + `$`)

/**
 * Skill Descript V3 parser
 * Started at: 2010-04-14
 *
 * Files:
 * - data/lua files/skillinfoz/skilldescript.lub
 *
 * Discontinued files:
 * - data/skilldesctable.txt
 * - data/skilldesctable2.txt
 *
 * Format:
 * SKILL_DESCRIPT = {
 * 	[SkillConst] = {
 * 		"Desc line 1",
 * 		"Desc line 2",
 * 		-- ...
 * 	}
 * }
 *
 * Notes:
 * - This aligns with Skill V4
 * - This update introduced auto-learning of pre-requirements, but I think this file itself
 *   was only changed for simplification purposes (replacing 2 txt for 1 lua)
 */
type SkillDescriptV3Parser struct {
}

// func NewSkillDescriptV3Parser() SkillDescriptParser {
// 	return &SkillDescriptV3Parser{}
// }

/* Not implemented because we are starting on 2012-01-01, which is pretty much using V4 */
