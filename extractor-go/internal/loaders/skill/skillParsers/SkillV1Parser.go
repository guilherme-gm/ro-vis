package skillParsers

/**
 * Skill V1 structure/parser
 * Started at: game launch (? 2001)
 *
 * Files:
 * - data/skilldesctable.txt (Description / skillDescript V1)
 * - data/skillnametable.txt (Part of Info / skillInfo V1)
 * - data/leveluseskillspamount.txt (Part of Info / skillInfo V1)
 */
type SkillV1Parser struct{}

func NewSkillV1Parser() *SkillV1Parser {
	return &SkillV1Parser{}
}

/* Not implemented due to extractor targeting 2012-01-01+ */
