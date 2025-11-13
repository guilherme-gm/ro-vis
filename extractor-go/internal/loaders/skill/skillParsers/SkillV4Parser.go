package skillParsers

/**
 * Skill V4 structure/parser
 * Started at: 2010-04-14
 *
 * Files:
 * - [NEW] data/lua files/skillinfoz/jobinheritlist.lub (Job Inherit V2, now uses consts)
 * - [NEW] data/lua files/skillinfoz/skillid.lub (Skill ID V1)
 * - [NEW] data/lua files/skillinfoz/skillinfolist.lub (Skill Info V2)
 * - [NEW] data/lua files/skillinfoz/skilldescript.lub (Skill Descript V3)
 * - [NEW] data/lua files/skillinfoz/skillinfo_f.lub (not parsed)
 * - [NEW] data/lua files/skillinfoz/skilltreeview.lub (not parsed)
 *
 * Discontinued files:
 * - data/skilldesctable.txt
 * - data/skilldesctable2.txt
 * - data/skillnametable.txt
 * - data/leveluseskillspamount.txt
 * - data/lua files/skillinfo/jobinheritlist.lub [Renamed/Small changes]
 * - data/lua files/skillinfo/skillinfo_f.lub [Renamed/Small changes]
 * - data/lua files/skillinfo/skilltreeview.lub [Renamed/Small changes]
 *
 * Notes:
 * - This update introduced auto-learning of pre-requirements
 */
type SkillV4Parser struct{}

func NewSkillV4Parser() *SkillV4Parser {
	return &SkillV4Parser{}
}

/* Not implemented due to extractor targeting 2012-01-01+ */
