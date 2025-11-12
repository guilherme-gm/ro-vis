package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type AddToHistoryResult struct {
	UpsertCount  int
	DeletedCount int
}

func (r AddToHistoryResult) String() string {
	return fmt.Sprintf("Upserted %d skills, deleted %d", r.UpsertCount, r.DeletedCount)
}

type SkillRepository struct {
	DB *database.Database
}

func NewSkillRepository(db *database.Database) *SkillRepository {
	return &SkillRepository{
		DB: db,
	}
}

func (r *SkillRepository) GetSkillJobs(tx *sql.Tx) ([]domain.SkillJob, error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.GetSkillJobs(context.Background())
	if err == sql.ErrNoRows {
		return []domain.SkillJob{}, nil
	}

	if err != nil {
		return nil, err
	}

	jobs := make([]domain.SkillJob, len(res))
	for idx, model := range res {
		jobs[idx] = model.ToDomain()
	}

	return jobs, nil
}

func (r *SkillRepository) insertSkillJob_sub(tx *sql.Tx, newJobs []*domain.SkillJob) error {
	queries := r.DB.GetDAO(tx)
	bulkParams := []dao.BulkInsertSkillJobParams{}
	for _, it := range newJobs {
		bulkParams = append(bulkParams, dao.BulkInsertSkillJobParams{
			Constant:      it.Constant,
			JobId:         it.JobId,
			InheritedJob:  sql.NullInt32(it.InheritedJob),
			InheritedJob2: sql.NullInt32(it.InheritedJob2),
			FirstUpdate:   it.FirstUpdate,
			LastUpdate:    it.LastUpdate,
			Deleted:       it.Deleted,
		})
	}

	_, err := queries.BulkInsertSkillJob(context.Background(), bulkParams)
	fmt.Println("inserted", len(bulkParams))
	return err
}

func (r *SkillRepository) AddSkillJobs(tx *sql.Tx, newJobs []*domain.SkillJob) error {
	if len(newJobs) == 0 {
		return nil
	}

	steps := (len(newJobs) / 500)

	i := 0
	for ; i < steps; i++ {
		slice := newJobs[i*500 : (i+1)*500]
		if err := r.insertSkillJob_sub(tx, slice); err != nil {
			return err
		}
	}

	slice := newJobs[i*500:]
	if err := r.insertSkillJob_sub(tx, slice); err != nil {
		return err
	}

	return nil
}

func (r *SkillRepository) GetCurrentSkills(tx *sql.Tx) ([]domain.Skill, error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.GetCurrentSkills(context.Background())
	if err == sql.ErrNoRows {
		return []domain.Skill{}, nil
	}

	if err != nil {
		return nil, err
	}

	skills := make([]domain.Skill, len(res))
	for idx, qmodel := range res {
		skills[idx] = qmodel.ToDomain()
	}

	return skills, nil
}

func (r *SkillRepository) insertSkill_sub(tx *sql.Tx, update string, newSkills []*domain.Skill) (AddToHistoryResult, error) {
	queries := r.DB.GetDAO(tx)

	// SkillID -> deleted (true | false)
	isWasDeleted := make(map[int32]bool, len(newSkills))
	bulkParams := []dao.BulkInsertSkillHistoryParams{}
	for _, it := range newSkills {
		isWasDeleted[it.SkillID] = it.Deleted
		spCostJson := domain.NewNullableNullString()
		if len(it.SpCost) > 0 {
			jsonBytes, _ := json.Marshal(it.SpCost)
			spCostJson = domain.NewNullableString(string(jsonBytes))
		}

		attackRangeJson := domain.NewNullableNullString()
		if len(it.AttackRange) > 0 {
			jsonBytes, _ := json.Marshal(it.AttackRange)
			attackRangeJson = domain.NewNullableString(string(jsonBytes))
		}

		needSkillListJson := domain.NewNullableNullString()
		if len(it.RequiredSkills) > 0 {
			jsonBytes, _ := json.Marshal(it.RequiredSkills)
			needSkillListJson = domain.NewNullableString(string(jsonBytes))
		}

		jobRequiredSkillsJson := domain.NewNullableNullString()
		if len(it.JobRequiredSkills) > 0 {
			jsonBytes, _ := json.Marshal(it.JobRequiredSkills)
			jobRequiredSkillsJson = domain.NewNullableString(string(jsonBytes))
		}

		insertParams := dao.BulkInsertSkillHistoryParams{
			PreviousHistoryID: sql.NullInt32(it.PreviousHistoryID),
			SkillId:           it.SkillID,
			FileVersion:       it.FileVersion,
			Update:            update,
		}

		if !it.Deleted {
			insertParams.Constant = sql.NullString(it.Constant)
			insertParams.Name = sql.NullString(it.Name)
			insertParams.Description = sql.NullString(it.Description)
			insertParams.MaxLevel = sql.NullInt32(it.MaxLevel)
			insertParams.SpCost = sql.NullString(spCostJson)
			insertParams.CanSelectLevel = sql.NullBool(it.CanSelectLevel)
			insertParams.AttackRange = sql.NullString(attackRangeJson)
			insertParams.RequiredSkills = sql.NullString(needSkillListJson)
			insertParams.JobRequiredSkills = sql.NullString(jobRequiredSkillsJson)
		}

		bulkParams = append(bulkParams, insertParams)
	}

	addResult := AddToHistoryResult{}
	_, err := queries.BulkInsertSkillHistory(context.Background(), bulkParams)
	if err != nil {
		return addResult, err
	}

	res, err := queries.GetSkillsIdsInUpdate(context.Background(), update)
	if err != nil {
		return addResult, err
	}

	upsertParams := []dao.BulkUpsertSkillParams{}
	for _, id := range res {
		if deleted, ok := isWasDeleted[id.SkillID]; ok {
			if deleted {
				addResult.DeletedCount++
			} else {
				addResult.UpsertCount++
			}

			upsertParams = append(upsertParams, dao.BulkUpsertSkillParams{
				SkillId:   id.SkillID,
				HistoryID: id.HistoryID,
				Deleted:   deleted,
			})
		}
	}

	_, err = queries.BulkUpsertSkills(context.Background(), upsertParams)
	if err != nil {
		return addResult, err
	}

	return addResult, nil
}

func (r *SkillRepository) AddSkillsToHistory(tx *sql.Tx, update string, newSkills []*domain.Skill) (AddToHistoryResult, error) {
	if len(newSkills) == 0 {
		return AddToHistoryResult{}, nil
	}

	finalResult := AddToHistoryResult{}
	steps := (len(newSkills) / 500)

	i := 0
	for ; i < steps; i++ {
		slice := newSkills[i*500 : (i+1)*500]

		res, err := r.insertSkill_sub(tx, update, slice)
		if err != nil {
			return AddToHistoryResult{}, err
		}
		finalResult.UpsertCount += res.UpsertCount
		finalResult.DeletedCount += res.DeletedCount
	}

	slice := newSkills[i*500:]

	res, err := r.insertSkill_sub(tx, update, slice)
	if err != nil {
		return AddToHistoryResult{}, err
	}
	finalResult.UpsertCount += res.UpsertCount
	finalResult.DeletedCount += res.DeletedCount

	return finalResult, nil
}
