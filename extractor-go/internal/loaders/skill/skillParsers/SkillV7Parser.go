package skillParsers

import (
	"fmt"
	"regexp"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/descriptParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/idParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/infoParser"
)

/**
 * Skill V7 structure/parser
 * Started at: 2020-09-15
 *
 * ApAmount added to SkillInfo.lub (SkillInfo V4)
 *
 * Files:
 * - data/luafiles514/lua files/skillinfoz/jobinheritlist.lub (Job Inherit V2, now uses consts)
 * - data/luafiles514/lua files/skillinfoz/skillid.lub (Skill ID V2)
 * - data/luafiles514/lua files/skillinfoz/skillinfolist.lub (Skill Info V4)
 * - data/luafiles514/lua files/skillinfoz/skilldescript.lub (Skill Descript V4)
 * - data/luafiles514/lua files/skillinfoz/skillinfo_f.lub (not parsed)
 * - data/luafiles514/lua files/skillinfoz/skilltreeview.lub (not parsed)
 */
type SkillV7Parser struct{}

func NewSkillV7Parser() *SkillV7Parser {
	return &SkillV7Parser{}
}

func (l *SkillV7Parser) FileVersion() int32 {
	return 7
}

func (l *SkillV7Parser) IsUpdateInRange(update *domain.Update) bool {
	return update.Date.After(time.Date(2020, time.September, 14, 0, 0, 0, 0, time.UTC)) &&
		update.Date.Before(time.Date(2025, time.June, 18, 0, 0, 0, 0, time.UTC))
}

func (l *SkillV7Parser) loadSkillIDs(basePath string, update domain.Update, skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill]) {
	change, err := update.GetChangeForFile(idParser.SkillIdV2)
	if err != nil {
		panic(err)
	}

	skillList := idParser.NewSkillIdV2Parser().Parse(basePath, &change)

	fileSkillMap := make(map[int32]bool)
	for fileConst, fileId := range skillList {
		if shouldSkipConst(fileConst) {
			fmt.Printf("Skipping SKID %s\n", fileConst)
			continue
		}

		if _, ok := fileSkillMap[int32(fileId)]; ok {
			panic(fmt.Sprintf("Duplicated SKID %d. 2nd constant found: %s", fileId, fileConst))
		}

		fileSkillMap[int32(fileId)] = true
		existingSkill, exists := skillUpdater.GetForRead(int32(fileId))
		shouldSave := false

		if !exists || existingSkill.SkillID != int32(fileId) || existingSkill.Constant.String != fileConst {
			shouldSave = true
		}

		if shouldSave {
			newSkill := skillUpdater.GetForEdit(int32(fileId))
			newSkill.Constant = domain.NewNullableString(fileConst)
			newSkill.SkillID = int32(fileId)
		}
	}

	for _, existingSkill := range skillUpdater.CurrentValues {
		if _, ok := fileSkillMap[existingSkill.SkillID]; !ok {
			deletedSkill := skillUpdater.GetForEdit(existingSkill.SkillID)
			deletedSkill.Deleted = true
		}
	}
}

func (l *SkillV7Parser) loadSkillInfos(
	basePath string,
	update domain.Update,
	skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill],
	jobUpdater *loaders.Updater[string, domain.SkillJob, *domain.SkillJob],
) {
	change, err := update.GetChangeForFile(infoParser.SkillInfoListV2)
	if err != nil {
		panic(err)
	}

	jobMap := make(map[string]int)
	jobUpdater.ForEach(func(key string, value domain.SkillJob) {
		jobMap[value.Constant] = int(value.JobId)
	})

	skillIdMap := make(map[string]int)
	skillUpdater.ForEach(func(key int32, value domain.Skill) {
		skillIdMap[value.Constant.String] = int(value.SkillID)
	})

	skillParser := infoParser.NewSkillInfoV4Parser() // change
	skillList := skillParser.Parse(basePath, &change, jobMap, skillIdMap)

	fileSkillMap := make(map[int32]bool)
	for _, fileSkill := range skillList {
		fileSkillMap[int32(fileSkill.SkillId)] = true
		existingSkill, exists := skillUpdater.GetForRead(int32(fileSkill.SkillId))
		shouldSave := false

		fileSkillDomain := fileSkill.ToSkill(existingSkill)
		if !exists || !existingSkill.Equals(fileSkillDomain) {
			shouldSave = true
		}

		if shouldSave {
			newSkill := skillUpdater.GetForEdit(int32(fileSkill.SkillId))
			*newSkill = fileSkill.ToSkill(*newSkill)
		}
	}

	for _, existingSkill := range skillUpdater.CurrentValues {
		if _, ok := fileSkillMap[existingSkill.SkillID]; !ok {
			if existingSkill.Description.Valid {
				deletedSkill := skillUpdater.GetForEdit(existingSkill.SkillID)
				deletedSkill.MaxLevel = domain.NewNullableNullInt32()
				deletedSkill.Name = domain.NewNullableNullString()
				deletedSkill.SpCost = nil
				deletedSkill.CanSelectLevel = domain.NewNullableNullBool()
				deletedSkill.AttackRange = nil
				deletedSkill.RequiredSkills = nil
				deletedSkill.JobRequiredSkills = nil
				deletedSkill.SkillScale = nil
				deletedSkill.ApCost = nil
			}
		}
	}
}

func (l *SkillV7Parser) loadSkillDesc(
	basePath string,
	update domain.Update,
	skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill],
) {
	change, err := update.GetChangeForFile(descriptParser.SkillDescriptV4)
	if err != nil {
		panic(err)
	}

	skillIdMap := make(map[string]int)
	skillUpdater.ForEach(func(key int32, value domain.Skill) {
		skillIdMap[value.Constant.String] = int(value.SkillID)
	})

	skillParser := descriptParser.NewSkillDescriptV4Parser()
	skillList := skillParser.Parse(basePath, &change, skillIdMap)

	fileSkillMap := make(map[int32]bool)
	for _, fileSkill := range skillList {
		fileSkillMap[int32(fileSkill.SkillId)] = true
		existingSkill, exists := skillUpdater.GetForRead(int32(fileSkill.SkillId))
		shouldSave := false

		fileSkillDomain := fileSkill.ToSkill(existingSkill)
		if !exists || !existingSkill.Equals(fileSkillDomain) {
			shouldSave = true
		}

		if shouldSave {
			newSkill := skillUpdater.GetForEdit(int32(fileSkill.SkillId))
			*newSkill = fileSkill.ToSkill(*newSkill)
		}
	}

	for _, existingSkill := range skillUpdater.CurrentValues {
		if _, ok := fileSkillMap[existingSkill.SkillID]; !ok {
			if existingSkill.Description.Valid {
				deletedSkill := skillUpdater.GetForEdit(existingSkill.SkillID)
				deletedSkill.Description = domain.NewNullableNullString()
			}
		}
	}
}

func (p *SkillV7Parser) Parse(basePath string, update domain.Update, skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill], jobUpdater *loaders.Updater[string, domain.SkillJob, *domain.SkillJob]) {
	if update.HasChangedAnyFiles([]*regexp.Regexp{idParser.SkillIdV2Regex}) {
		p.loadSkillIDs(basePath, update, skillUpdater)
	}
	if update.HasChangedAnyFiles([]*regexp.Regexp{infoParser.SkillInfoListV3Regex}) {
		p.loadSkillInfos(basePath, update, skillUpdater, jobUpdater)
	}
	if update.HasChangedAnyFiles([]*regexp.Regexp{descriptParser.SkillDescriptV4Regex}) {
		p.loadSkillDesc(basePath, update, skillUpdater)
	}
}
