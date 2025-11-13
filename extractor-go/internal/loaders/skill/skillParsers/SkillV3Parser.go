package skillParsers

/**
 * Skill V3 structure/parser
 * Started at: 2009-10-07
 *
 * Files:
 * - data/skilldesctable.txt (Description / skillDescript V2 -- same as V1)
 * - data/skilldesctable2.txt (Description / skillDescript V2)
 * - data/skillnametable.txt (Part of Info / skillInfo V1)
 * - data/leveluseskillspamount.txt (Part of Info / skillInfo V1)
 * - [NEW] data/lua files/skillinfo/jobinheritlist.lub (Job Inherit V1)
 * - [NEW] data/lua files/skillinfo/skillinfo_f.lub (not parsed)
 * - [NEW] data/lua files/skillinfo/skilltreeview.lub (not parsed)
 *
 * Discontinued files:
 * - data/skilltreeview.txt
 *
 * Notes:
 * It is unclear whether it brought any visual changes or were simply part of internal refactoring.
 * The patch notes around the date doesn't mention it.
 */
type SkillV3Parser struct{}

func NewSkillV3Parser() *SkillV3Parser {
	return &SkillV3Parser{}
}

/* Not implemented due to extractor targeting 2012-01-01+ */
