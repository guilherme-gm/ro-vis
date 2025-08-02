package dao

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func (q *GetCurrentI18nsRow) ToDomain() domain.I18n {
	return domain.I18n{
		PreviousHistoryID: domain.NullableInt64(q.PreviousHistoryID),
		HistoryID:         ToNullableInt64(q.HistoryID),
		I18nId:            q.I18nID,
		FileVersion:       q.FileVersion,
		ContainerFile:     q.ContainerFile,
		EnText:            q.EnText,
		PtBrText:          q.PtBrText,
		Active:            q.Active,
		Deleted:           q.Deleted,
	}
}

func (q *I18nHistory) ToDomain() domain.I18n {
	return domain.I18n{
		PreviousHistoryID: domain.NullableInt64(q.PreviousHistoryID),
		HistoryID:         ToNullableInt64(q.HistoryID),
		I18nId:            q.I18nID,
		FileVersion:       q.FileVersion,
		ContainerFile:     q.ContainerFile,
		EnText:            q.EnText,
		PtBrText:          q.PtBrText,
		Active:            q.Active,
	}
}

func (q *PreviousI18nHistoryVw) ToDomain() domain.I18n {
	return domain.I18n{
		PreviousHistoryID: domain.NullableInt64(q.PreviousHistoryID),
		HistoryID:         domain.NullableInt64(q.HistoryID),
		I18nId:            q.I18nID.Val,
		FileVersion:       q.FileVersion.Int32,
		ContainerFile:     q.ContainerFile.String,
		EnText:            q.EnText.String,
		PtBrText:          q.PtBrText.String,
		Active:            q.Active.Bool,
	}
}
