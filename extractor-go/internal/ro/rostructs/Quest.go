package rostructs

import (
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type QuestV1 struct {
	QuestId     int32 `lua:"@index"`
	Title       string
	Description []string
	Summary     string
}

func (q *QuestV1) ToDomain() domain.Quest {
	return domain.Quest{
		QuestID:     int32(q.QuestId),
		FileVersion: 1,
		Title:       dao.ToNullString(q.Title),
		Summary:     dao.ToNullString(q.Summary),
		Description: dao.ToNullString(strings.Join(q.Description, "\n")),
	}
}

type QuestV3RewardItem struct {
	ItemID  int32
	ItemNum int32
}

type QuestV3 struct {
	QuestId        int32 `lua:"@index"`
	Title          string
	Description    []string
	Summary        string
	IconName       string
	NpcSpr         string
	NpcNavi        string
	NpcPosX        int32
	NpcPosY        int32
	RewardEXP      string
	RewardJEXP     string
	RewardItemList []QuestV3RewardItem
}

func (q *QuestV3) ToDomain() domain.Quest {
	return domain.Quest{
		QuestID:     int32(q.QuestId),
		FileVersion: 3,
		Title:       dao.ToNullString(q.Title),
		Description: dao.ToNullString(strings.Join(q.Description, "\n")),
		Summary:     dao.ToNullString(q.Summary),
		IconName:    dao.ToNullString(q.IconName),
		NpcSpr:      dao.ToNullString(q.NpcSpr),
		NpcNavi:     dao.ToNullString(q.NpcNavi),
		NpcPosX:     dao.ToNullInt32(q.NpcPosX),
		NpcPosY:     dao.ToNullInt32(q.NpcPosY),
		RewardExp:   dao.ToNullString(q.RewardEXP),
		RewardJexp:  dao.ToNullString(q.RewardJEXP),
		// TODO:
		// RewardItemList: dao.ToNullString(q.RewardItemList),
	}
}

type QuestV4 struct {
	QuestId        int32 `lua:"@index"`
	Title          string
	Description    []string
	Summary        string
	IconName       string
	NpcSpr         string
	NpcNavi        string
	NpcPosX        int32
	NpcPosY        int32
	RewardEXP      string
	RewardJEXP     string
	RewardItemList []QuestV4RewardItem
	// New in V4
	CoolTimeQuest int32
}

type QuestV4RewardItem struct {
	ItemID  int32
	ItemNum int32
}

func (q *QuestV4) ToDomain() domain.Quest {
	return domain.Quest{
		QuestID:     int32(q.QuestId),
		FileVersion: 3,
		Title:       dao.ToNullString(q.Title),
		Description: dao.ToNullString(strings.Join(q.Description, "\n")),
		Summary:     dao.ToNullString(q.Summary),
		IconName:    dao.ToNullString(q.IconName),
		NpcSpr:      dao.ToNullString(q.NpcSpr),
		NpcNavi:     dao.ToNullString(q.NpcNavi),
		NpcPosX:     dao.ToNullInt32(q.NpcPosX),
		NpcPosY:     dao.ToNullInt32(q.NpcPosY),
		RewardExp:   dao.ToNullString(q.RewardEXP),
		RewardJexp:  dao.ToNullString(q.RewardJEXP),
		// TODO:
		// RewardItemList: dao.ToNullString(q.RewardItemList),
	}
}
