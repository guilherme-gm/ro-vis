package skillParsers

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/descriptParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/idParser"
	"github.com/guilherme-gm/ro-vis/extractor/internal/loaders/skill/skillParsers/infoParser"
)

/**
 * Skill V5 starts on 2012-01-11, and it is a simple directory change compared to V4.
 * Files where moved from data/lua files/skillinfoz/* to data/luafiles514/lua files/skillinfoz/*
 */
type SkillV5Parser struct{}

func NewSkillV5Parser() *SkillV5Parser {
	return &SkillV5Parser{}
}

func (l *SkillV5Parser) loadSkillIDs(basePath string, update domain.Update, skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill]) {
	change, err := update.GetChangeForFile(idParser.SkillIdV2)
	if err != nil {
		panic(err)
	}

	skillList := idParser.NewSkillIdV2Parser().Parse(basePath, &change)

	fileSkillMap := make(map[int32]bool)
	for fileConst, fileId := range skillList {
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
			*newSkill = fileSkillDomain
		}
	}

	for _, existingSkill := range skillUpdater.CurrentValues {
		if _, ok := fileSkillMap[existingSkill.SkillID]; !ok {
			deletedSkill := skillUpdater.GetForEdit(existingSkill.SkillID)
			deletedSkill.Deleted = true
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
			*newSkill = fileSkillDomain
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

func (p *SkillV5Parser) Parse(basePath string, update domain.Update, skillUpdater *loaders.Updater[int32, domain.Skill, *domain.Skill], jobUpdater *loaders.Updater[string, domain.SkillJob, *domain.SkillJob]) {
	p.loadSkillIDs(basePath, update, skillUpdater)
	p.loadSkillInfos(basePath, update, skillUpdater, jobUpdater)
	p.loadSkillDesc(basePath, update, skillUpdater)
}
