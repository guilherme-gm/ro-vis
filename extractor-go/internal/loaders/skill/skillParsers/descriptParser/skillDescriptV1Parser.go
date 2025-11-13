package descriptParser

import (
	"regexp"
)

var SkillDescriptV1 = "data/skilldesctable.txt"
var SkillDescriptV1Regex = regexp.MustCompile(`(?i)^` + SkillDescriptV1 + `$`)

/**
 * Skill Descript V1 parser
 * Started at: game launch (?) (2001)
 *
 * Files:
 * - data/skilldesctable.txt
 *
 * Format:
 *  SKILL_CONST1#
 *  DESC
 *  DESC
 *  #
 *  SKILL_CONST2#
 *  DESC
 *  DESC
 *  #
 *
 * Notes:
 * - This aligns with Skill V1
 */
type SkillDescriptV1Parser struct {
}

// func NewSkillDescriptV1Parser() SkillDescriptParser {
// 	return &SkillDescriptV1Parser{}
// }

/* Not implemented because we are starting on 2012-01-01, which is using V4 */
