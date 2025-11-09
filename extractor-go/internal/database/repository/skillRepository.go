package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

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
