package skillParsers

/**
 * Skill V2 structure/parser
 * Started at: 2008-06-25
 *
 * Files:
 * - data/skilldesctable.txt (Description / skillDescript V2 -- same as V1)
 * - [NEW] data/skilldesctable2.txt (Description / skillDescript V2)
 * - data/skillnametable.txt (Part of Info / skillInfo V1)
 * - data/leveluseskillspamount.txt (Part of Info / skillInfo V1)
 * - [NEW] data/skilltreeview.txt (Currently not parsed)
 *
 * Notes:
 * - This change was made to support 2 skill window view (list vs tree)
 *    - skilldesctable for the minimized window
 *    - skilldesctable2 for the expanded window
 *    - pre-required skills are shown in the expanded view (skilldesctable2)
 */
type SkillV2Parser struct{}

func NewSkillV2Parser() *SkillV2Parser {
	return &SkillV2Parser{}
}

/* Not implemented due to extractor targeting 2012-01-01+ */
