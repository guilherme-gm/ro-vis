package repository

import (
	"context"
	"database/sql"
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
	HistoryBaseRepository[
		domain.Skill,
		dao.BulkInsertSkillHistoryParams,
		*dao.BulkInsertSkillHistoryParams,
		*dao.BulkUpsertSkillParams,
	]
}

func NewSkillRepository(db database.IDatabase) *SkillRepository {
	repo := &SkillRepository{
		HistoryBaseRepository: HistoryBaseRepository[
			domain.Skill,
			dao.BulkInsertSkillHistoryParams,
			*dao.BulkInsertSkillHistoryParams,
			*dao.BulkUpsertSkillParams,
		]{
			DB:       db,
			BulkSize: 500,
		},
	}
	repo.getCurrentData = repo.getCurrentDataFn
	repo.bulkInsertHistory = repo.bulkInsertSkillHistoryFn
	repo.bulkUpsertRecords = repo.bulkUpsertSkillFn
	repo.getIdsInUpdate = repo.getIdsInUpdateFn
	repo.newBulkParamEntry = func() *dao.BulkInsertSkillHistoryParams { return &dao.BulkInsertSkillHistoryParams{} }
	repo.newRecordParam = func() *dao.BulkUpsertSkillParams { return &dao.BulkUpsertSkillParams{} }
	return repo
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

func (r *SkillRepository) insertSkillJob_sub(tx *sql.Tx, newJobs []domain.SkillJob) error {
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

func (r *SkillRepository) AddSkillJobs(tx *sql.Tx, newJobs []domain.SkillJob) error {
	if len(newJobs) == 0 {
		return nil
	}

	steps := (len(newJobs) / r.BulkSize)

	i := 0
	for ; i < steps; i++ {
		slice := newJobs[i*r.BulkSize : (i+1)*r.BulkSize]
		if err := r.insertSkillJob_sub(tx, slice); err != nil {
			return err
		}
	}

	slice := newJobs[i*r.BulkSize:]
	if err := r.insertSkillJob_sub(tx, slice); err != nil {
		return err
	}

	return nil
}

func (r *SkillRepository) getCurrentDataFn(dao dao.IDAO, ctx context.Context) ([]domain.Skill, error) {
	return toDomainSlice(dao.GetCurrentSkills(ctx))
}

func (r *SkillRepository) bulkInsertSkillHistoryFn(dao dao.IDAO, arg []*dao.BulkInsertSkillHistoryParams) (sql.Result, error) {
	return dao.BulkInsertSkillHistory(context.Background(), arg)
}

func (r *SkillRepository) bulkUpsertSkillFn(dao dao.IDAO, arg []*dao.BulkUpsertSkillParams) (sql.Result, error) {
	return dao.BulkUpsertSkills(context.Background(), arg)
}

func (r *SkillRepository) getIdsInUpdateFn(dao dao.IDAO, update string) ([]IdHistory, error) {
	return toIdHistorySlice(dao.GetSkillsIdsInUpdate(context.Background(), update))
}

// func (r *SkillRepository) insertSkill_sub(tx *sql.Tx, update string, newSkills []*domain.Skill) (AddToHistoryResult, error) {
// 	queries := r.DB.GetDAO(tx)

// 	// SkillID -> deleted (true | false)
// 	isWasDeleted := make(map[int32]bool, len(newSkills))
// 	bulkParams := make([]dao.BulkInsertSkillHistoryParams, len(newSkills))
// 	for idx, it := range newSkills {
// 		isWasDeleted[it.SkillID] = it.Deleted
// 		bulkParams[idx].FillFromDomain(it, update)
// 	}

// 	addResult := AddToHistoryResult{}
// 	_, err := queries.BulkInsertSkillHistory(context.Background(), bulkParams)
// 	if err != nil {
// 		return addResult, err
// 	}

// 	res, err := queries.GetSkillsIdsInUpdate(context.Background(), update)
// 	if err != nil {
// 		return addResult, err
// 	}

// 	upsertParams := make([]dao.BulkUpsertSkillParams, len(res))
// 	for idx, id := range res {
// 		if deleted, ok := isWasDeleted[id.SkillID]; ok {
// 			if deleted {
// 				addResult.DeletedCount++
// 			} else {
// 				addResult.UpsertCount++
// 			}

// 			upsertParams[idx].Fill(id.SkillID, id.HistoryID, deleted)
// 		}
// 	}

// 	_, err = queries.BulkUpsertSkills(context.Background(), upsertParams)
// 	if err != nil {
// 		return addResult, err
// 	}

// 	return addResult, nil
// }

// func (r *SkillRepository) AddSkillsToHistory(tx *sql.Tx, update string, newSkills []*domain.Skill) (AddToHistoryResult, error) {
// 	if len(newSkills) == 0 {
// 		return AddToHistoryResult{}, nil
// 	}

// 	finalResult := AddToHistoryResult{}
// 	steps := (len(newSkills) / r.BulkSize)

// 	i := 0
// 	for ; i < steps; i++ {
// 		slice := newSkills[i*r.BulkSize : (i+1)*r.BulkSize]

// 		res, err := r.insertSkill_sub(tx, update, slice)
// 		if err != nil {
// 			return AddToHistoryResult{}, err
// 		}
// 		finalResult.UpsertCount += res.UpsertCount
// 		finalResult.DeletedCount += res.DeletedCount
// 	}

// 	slice := newSkills[i*r.BulkSize:]

// 	res, err := r.insertSkill_sub(tx, update, slice)
// 	if err != nil {
// 		return AddToHistoryResult{}, err
// 	}
// 	finalResult.UpsertCount += res.UpsertCount
// 	finalResult.DeletedCount += res.DeletedCount

// 	return finalResult, nil
// }

func (r *SkillRepository) CountSkills(tx *sql.Tx) (int32, error) {
	queries := r.DB.GetDAO(tx)

	res, err := queries.CountSkills(context.Background())
	if err == sql.ErrNoRows {
		return int32(res), nil
	}

	if err != nil {
		return 0, err
	}

	return int32(res), nil
}

func (r *SkillRepository) GetSkills(tx *sql.Tx, pagination Pagination) ([]domain.MinSkill, error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.GetSkillList(context.Background(), dao.GetSkillListParams{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []domain.MinSkill{}, nil
	}
	if err != nil {
		return []domain.MinSkill{}, nil
	}

	skills := make([]domain.MinSkill, len(res))
	for idx, val := range res {
		skills[idx] = domain.MinSkill{
			SkillID:    val.SkillID,
			LastUpdate: val.Lastupdate,
			Constant:   domain.NullableString(val.Constant),
			Name:       domain.NullableString(val.Name),
		}
	}
	return skills, nil
}

func (r *SkillRepository) sqlRecordToDomain(dbFrom dao.PreviousSkillHistoryVw, dbTo dao.SkillsHistory, lastUpdate sql.NullString) FromToRecord[domain.Skill] {
	var from *domain.Record[domain.Skill] = nil
	var to *domain.Record[domain.Skill] = nil

	if dbFrom.HistoryID.Valid {
		from = &domain.Record[domain.Skill]{
			Update: dbFrom.Update.String,
			Data:   dbFrom.ToDomain(),
		}
	}

	if dbTo.HistoryID != 0 {
		to = &domain.Record[domain.Skill]{
			Update: dbTo.Update,
			Data:   dbTo.ToDomain(),
		}
	}

	return FromToRecord[domain.Skill]{
		LastUpdate: domain.NullableString(lastUpdate),
		From:       from,
		To:         to,
	}
}

func (r *SkillRepository) GetSkillHistory(tx *sql.Tx, skillId int32, pagination Pagination) ([]FromToRecord[domain.Skill], error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.GetSkillHistory(context.Background(), dao.GetSkillHistoryParams{
		SkillID: skillId,
		Offset:  pagination.Offset,
		Limit:   pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []FromToRecord[domain.Skill]{}, nil
	}

	if err != nil {
		return nil, err
	}

	records := make([]FromToRecord[domain.Skill], len(res))
	for idx, qmodel := range res {
		records[idx] = r.sqlRecordToDomain(qmodel.PreviousSkillHistoryVw, qmodel.SkillsHistory, sql.NullString{})
	}

	return records, nil
}

func (r *SkillRepository) CountChangesInUpdate(tx *sql.Tx, update string) (int, error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.CountChangedSkillsInUpdate(context.Background(), update)
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *SkillRepository) GetChangesInUpdate(tx *sql.Tx, update string, pagination Pagination) ([]FromToRecord[domain.Skill], error) {
	queries := r.DB.GetDAO(tx)
	res, err := queries.GetChangedSkills(context.Background(), dao.GetChangedSkillsParams{
		Update: update,
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	})
	if err == sql.ErrNoRows {
		return []FromToRecord[domain.Skill]{}, nil
	}

	if err != nil {
		return nil, err
	}

	records := make([]FromToRecord[domain.Skill], len(res))
	for idx, qmodel := range res {
		records[idx] = r.sqlRecordToDomain(qmodel.PreviousSkillHistoryVw, qmodel.SkillsHistory, qmodel.Lastupdate)
	}

	return records, nil
}
