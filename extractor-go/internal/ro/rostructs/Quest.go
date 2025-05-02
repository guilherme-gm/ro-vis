package rostructs

import (
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
)

type QuestV1 struct {
	QuestId     int32 `lua:"@index"`
	Title       string
	Description []string
	Summary     string
}

func (q *QuestV1) ToModel() dao.QuestHistory {
	return dao.QuestHistory{
		QuestID:     int32(q.QuestId),
		FileVersion: 1,
		Title:       dao.ToNullString(q.Title),
		Summary:     dao.ToNullString(q.Summary),
		Description: dao.ToNullString(strings.Join(q.Description, "\n")),
	}
}
