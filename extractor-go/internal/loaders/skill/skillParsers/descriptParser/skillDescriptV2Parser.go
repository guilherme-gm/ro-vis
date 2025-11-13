package descriptParser

import (
	"regexp"
)

var SkillDescriptV2 = "data/skilldesctable.txt"
var SkillDescriptV2Regex = regexp.MustCompile(`(?i)^` + SkillDescriptV2 + `$`)

var SkillDescript2V2 = "data/skilldesctable2.txt"
var SkillDescript2V2Regex = regexp.MustCompile(`(?i)^` + SkillDescript2V2 + `$`)

/**
 * Skill Descript V2 parser
 * Started at: 2008-06-25
 *
 * Files:
 * - data/skilldesctable.txt
 * - data/skilldesctable2.txt
 *
 * Format: (for both files)
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
 * - This aligns with Skill V2
 * - This change was made to support 2 skill window view, skilldesctable for the minimized window
 *   and skilldesctable2 for the expanded window, which shows the pre-required skills (in the description).
 */
type SkillDescriptV2Parser struct {
}

// func NewSkillDescriptV2Parser() SkillDescriptParser {
// 	return &SkillDescriptV2Parser{}
// }

/* Not implemented because we are starting on 2012-01-01, which is using V4 */
