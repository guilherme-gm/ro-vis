package jobParsers

import (
	"regexp"
)

var JobInheritListV2 = "data/lua files/skillinfoz/jobinheritlist.lub"
var JobInheritListV2Regex = regexp.MustCompile(`(?i)^` + JobInheritListV2 + `$`)

/**
 * V2 starts on 2010-04-14, it replaces the previous lua files/skillinfo/jobinheritlist.lub
 * Changes:
 * - Introduced JOBID constants ( JobConst = JobId )
 * - JOB_INHERIT_LIST / JOB_INHERIT_LIST2 now uses constants
 *
 * Notes:
 * - This date may be linked to Sakray update
 * - It aligns with the addition of SkillInfoList
 *    - And seems to have had the immediate effect of allowing automated point allocation
 *      for pre-requisites when you click to put points in a new skill
 * - It aligns with Skill V4
 */
type JobInheritListV2Parser struct {
}

/* Not implemented because we are starting on 2012-01-01, which is pretty much at V3 */
