package skillParsers

import (
	"fmt"
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/descriptParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/idParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/infoParser"
)

/**
 * Skill V5 structure/parser
 * Started at: 2012-01-11
 * Simply moved files from data/lua files/skillinfoz/* to data/lua files514/lua files/skillinfoz/*
 *
 * Files:
 * - data/luafiles514/lua files/skillinfoz/jobinheritlist.lub (Job Inherit V2, now uses consts)
 * - data/luafiles514/lua files/skillinfoz/skillid.lub (Skill ID V2)
 * - data/luafiles514/lua files/skillinfoz/skillinfolist.lub (Skill Info V3)
 * - data/luafiles514/lua files/skillinfoz/skilldescript.lub (Skill Descript V4)
 * - data/luafiles514/lua files/skillinfoz/skillinfo_f.lub (not parsed)
 * - data/luafiles514/lua files/skillinfoz/skilltreeview.lub (not parsed)
 */
type SkillV5Parser struct{}

func NewSkillV5Parser() *SkillV5Parser {
	return &SkillV5Parser{}
}

func shouldSkipConst(constant string) bool {
	// These were present in a few 2015 updates of SkillID,
	// but they were not used and they conflict with other real skills
	// They were later removed in 2015-12.
	// They cause a lot of issues since we base on the numeric ID as the lead, and
	// no reason to fret over those 4, so they are skipped.
	return constant == "SR_FLASHCOMBO_ATK_STEP1" ||
		constant == "SR_FLASHCOMBO_ATK_STEP2" ||
		constant == "SR_FLASHCOMBO_ATK_STEP3" ||
		constant == "SR_FLASHCOMBO_ATK_STEP4"
}

func (l *SkillV5Parser) loadSkillIDs(basePath string, update domain.Update, skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill]) {
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

		fileSkillMap[int32(fileId)] = true
		existingSkill, exists := skillUpdater.GetForRead(int32(fileId))
		shouldSave := false

		if !exists || existingSkill.SkillID != int32(fileId) || existingSkill.Constant.String != fileConst {
			shouldSave = true
		}

		if shouldSave {
			newSkill := skillUpdater.GetForEdit(int32(fileId))
			newSkill.FileVersion = 5
			newSkill.Constant = domain.NewNullableString(fileConst)
			newSkill.SkillID = int32(fileId)
		}
	}

	for _, existingSkill := range skillUpdater.CurrentValues {
		if _, ok := fileSkillMap[existingSkill.SkillID]; !ok {
			deletedSkill := skillUpdater.GetForEdit(existingSkill.SkillID)
			deletedSkill.FileVersion = 5
			deletedSkill.Deleted = true
		}
	}
}

type KV[T comparable] struct {
	Key   T
	Value int
}

func (l *SkillV5Parser) loadSkillInfos(
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

	skillParser := infoParser.NewSkillInfoV2Parser()
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
			newSkill.FileVersion = 5
			*newSkill = fileSkillDomain
		}
	}
}

func (l *SkillV5Parser) loadSkillDesc(
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
			newSkill.FileVersion = 5
			*newSkill = fileSkillDomain
		}
	}

	for _, existingSkill := range skillUpdater.CurrentValues {
		if _, ok := fileSkillMap[existingSkill.SkillID]; !ok {
			if existingSkill.Description.Valid {
				deletedSkill := skillUpdater.GetForEdit(existingSkill.SkillID)
				deletedSkill.FileVersion = 5
				deletedSkill.Description = domain.NewNullableNullString()
			}
		}
	}
}

func (p *SkillV5Parser) Parse(basePath string, update domain.Update, skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill], jobUpdater *loaders.Updater[string, domain.SkillJob, *domain.SkillJob]) {
	if update.HasChangedAnyFiles([]*regexp.Regexp{idParser.SkillIdV2Regex}) {
		p.loadSkillIDs(basePath, update, skillUpdater)
	}
	if update.HasChangedAnyFiles([]*regexp.Regexp{infoParser.SkillInfoListV2Regex}) {
		p.loadSkillInfos(basePath, update, skillUpdater, jobUpdater)
	}
	if update.HasChangedAnyFiles([]*regexp.Regexp{descriptParser.SkillDescriptV4Regex}) {
		p.loadSkillDesc(basePath, update, skillUpdater)
	}
}
