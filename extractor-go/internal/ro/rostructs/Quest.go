package rostructs

import (
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/dao"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type QuestV1 struct {
	QuestId     int
	Title       string
	Icon        string
	Image       string
	Description string
	Summary     string
}

func (q *QuestV1) ToDomain() domain.Quest {
	return domain.Quest{
		QuestID:     int32(q.QuestId),
		FileVersion: 1,
		Title:       dao.ToNullString(q.Title),
		IconName:    dao.ToNullString(q.Icon),
		OldImage:    dao.ToNullString(q.Image),
		Summary:     dao.ToNullString(q.Summary),
		Description: dao.ToNullString(q.Description),
	}
}

type QuestV3RewardItem struct {
	ItemID  int
	ItemNum int
}

type QuestV3 struct {
	QuestId        int `lua:"@index"`
	Title          string
	Description    []string
	Summary        string
	IconName       string
	NpcSpr         string
	NpcNavi        string
	NpcPosX        int
	NpcPosY        int
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
		NpcPosX:     dao.ToNullInt32(int32(q.NpcPosX)),
		NpcPosY:     dao.ToNullInt32(int32(q.NpcPosY)),
		RewardExp:   dao.ToNullString(q.RewardEXP),
		RewardJexp:  dao.ToNullString(q.RewardJEXP),
		// TODO:
		// RewardItemList: dao.ToNullString(q.RewardItemList),
	}
}

type QuestV4 struct {
	QuestId        int `lua:"@index"`
	Title          string
	Description    []string
	Summary        string
	IconName       string
	NpcSpr         string
	NpcNavi        string
	NpcPosX        int
	NpcPosY        int
	RewardEXP      string
	RewardJEXP     string
	RewardItemList []QuestV4RewardItem
	// New in V4
	CoolTimeQuest int
}

type QuestV4RewardItem struct {
	ItemID  int
	ItemNum int
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
		NpcPosX:     dao.ToNullInt32(int32(q.NpcPosX)),
		NpcPosY:     dao.ToNullInt32(int32(q.NpcPosY)),
		RewardExp:   dao.ToNullString(q.RewardEXP),
		RewardJexp:  dao.ToNullString(q.RewardJEXP),
		// TODO:
		// RewardItemList: dao.ToNullString(q.RewardItemList),
	}
}
