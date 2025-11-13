package idParser

import (
	"regexp"
)

var SkillIdV1 = "data/luafiles514/lua files/skillinfoz/skillid.lub"
var SkillIdV1Regex = regexp.MustCompile(`(?i)^` + SkillIdV1 + `$`)

/**
 * Skill ID V1 parser
 * Started at: 2010-04-14
 *
 * Files:
 * - data/lua files/skillinfoz/skillid.lub
 *
 * Format:
 * SKID = { SkillConst = SkillId }
 *
 * Notes:
 * - This aligns with Skill V4
 * - This update introduced auto-learning of pre-requirements, but I think this file itself
 *   was only changed for simplification purposes (having proper Lua constants)
 */
type SkillIdV1Parser struct {
}
