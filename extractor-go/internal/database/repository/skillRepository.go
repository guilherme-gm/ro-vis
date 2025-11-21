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
	BaseRepository
}

func NewSkillRepository(db database.IDatabase) *SkillRepository {
	return &SkillRepository{
		BaseRepository: BaseRepository{DB: db},
	}
}

func (r *SkillRepository) GetSkillJobs(tx *sql.Tx) ([]domain.SkillJob, error) {
	queries := r.dao(tx)
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

	for _, seg := range chunkIndices(len(newJobs), 500) {
		slice := newJobs[seg[0]:seg[1]]
		if err := r.insertSkillJob_sub(tx, slice); err != nil {
			return err
		}
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

		apCostJson := domain.NewNullableNullString()
		if len(it.ApCost) > 0 {
			jsonBytes, _ := json.Marshal(it.ApCost)
			apCostJson = domain.NewNullableString(string(jsonBytes))
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

		skillScaleJson := domain.NewNullableNullString()
		if len(it.SkillScale) > 0 {
			jsonBytes, _ := json.Marshal(it.SkillScale)
			skillScaleJson = domain.NewNullableString(string(jsonBytes))
		}

		castFlagsJson := domain.NewNullableNullString()
		if len(it.CastFlags) > 0 {
			jsonBytes, _ := json.Marshal(it.CastFlags)
			castFlagsJson = domain.NewNullableString(string(jsonBytes))
		}

		castFixedDelayJson := domain.NewNullableNullString()
		if len(it.CastFixedDelay) > 0 {
			jsonBytes, _ := json.Marshal(it.CastFixedDelay)
			castFixedDelayJson = domain.NewNullableString(string(jsonBytes))
		}

		castStatDelayJson := domain.NewNullableNullString()
		if len(it.CastStatDelay) > 0 {
			jsonBytes, _ := json.Marshal(it.CastStatDelay)
			castStatDelayJson = domain.NewNullableString(string(jsonBytes))
		}

		singlePostDelayJson := domain.NewNullableNullString()
		if len(it.SinglePostDelay) > 0 {
			jsonBytes, _ := json.Marshal(it.SinglePostDelay)
			singlePostDelayJson = domain.NewNullableString(string(jsonBytes))
		}

		globalPostDelayJson := domain.NewNullableNullString()
		if len(it.GlobalPostDelay) > 0 {
			jsonBytes, _ := json.Marshal(it.GlobalPostDelay)
			globalPostDelayJson = domain.NewNullableString(string(jsonBytes))
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
			insertParams.IsPassive = sql.NullBool(it.IsPassive)
			insertParams.SpCost = sql.NullString(spCostJson)
			insertParams.ApCost = sql.NullString(apCostJson)
			insertParams.CanSelectLevel = sql.NullBool(it.CanSelectLevel)
			insertParams.AttackRange = sql.NullString(attackRangeJson)
			insertParams.RequiredSkills = sql.NullString(needSkillListJson)
			insertParams.JobRequiredSkills = sql.NullString(jobRequiredSkillsJson)
			insertParams.SkillScale = sql.NullString(skillScaleJson)
			insertParams.CastFlags = sql.NullString(castFlagsJson)
			insertParams.CastFixedDelay = sql.NullString(castFixedDelayJson)
			insertParams.CastStatDelay = sql.NullString(castStatDelayJson)
			insertParams.SinglePostDelay = sql.NullString(singlePostDelayJson)
			insertParams.GlobalPostDelay = sql.NullString(globalPostDelayJson)
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
	for _, seg := range chunkIndices(len(newSkills), 500) {
		slice := newSkills[seg[0]:seg[1]]
		res, err := r.insertSkill_sub(tx, update, slice)
		if err != nil {
			return AddToHistoryResult{}, err
		}
		finalResult.UpsertCount += res.UpsertCount
		finalResult.DeletedCount += res.DeletedCount
	}

	return finalResult, nil
}

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
