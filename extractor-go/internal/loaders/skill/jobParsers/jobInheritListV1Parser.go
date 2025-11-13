package jobParsers

import (
	"regexp"
)

var JobInheritListV1 = "data/lua files/skillinfoz/jobinheritlist.lub"
var JobInheritListV1Regex = regexp.MustCompile(`(?i)^` + JobInheritListV1 + `$`)

/**
 * V1 starts on 2009-10-07, it is the first version doing it
 * It is probably linked with the SkillTree conversion to Lua.
 *
 * It is unclear whether it brought any visual changes or were simply part of internal refactoring.
 * The patch notes around the date doesn't mention it.
 *
 * It aligns with Skill V3
 */
type JobInheritListV1Parser struct {
}

/* Not implemented because we are starting on 2012-01-01, which is pretty much at V3 */
