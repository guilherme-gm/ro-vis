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
		Title:       dao.ToNullableString(q.Title),
		IconName:    dao.ToNullableString(q.Icon),
		OldImage:    dao.ToNullableString(q.Image),
		Summary:     dao.ToNullableString(q.Summary),
		Description: dao.ToNullableString(q.Description),
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
		Title:       dao.ToNullableString(q.Title),
		Description: dao.ToNullableString(strings.Join(q.Description, "\n")),
		Summary:     dao.ToNullableString(q.Summary),
		IconName:    dao.ToNullableString(q.IconName),
		NpcSpr:      dao.ToNullableString(q.NpcSpr),
		NpcNavi:     dao.ToNullableString(q.NpcNavi),
		NpcPosX:     dao.ToNullableInt32(int32(q.NpcPosX)),
		NpcPosY:     dao.ToNullableInt32(int32(q.NpcPosY)),
		RewardEXP:   dao.ToNullableString(q.RewardEXP),
		RewardJEXP:  dao.ToNullableString(q.RewardJEXP),
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
		Title:       dao.ToNullableString(q.Title),
		Description: dao.ToNullableString(strings.Join(q.Description, "\n")),
		Summary:     dao.ToNullableString(q.Summary),
		IconName:    dao.ToNullableString(q.IconName),
		NpcSpr:      dao.ToNullableString(q.NpcSpr),
		NpcNavi:     dao.ToNullableString(q.NpcNavi),
		NpcPosX:     dao.ToNullableInt32(int32(q.NpcPosX)),
		NpcPosY:     dao.ToNullableInt32(int32(q.NpcPosY)),
		RewardEXP:   dao.ToNullableString(q.RewardEXP),
		RewardJEXP:  dao.ToNullableString(q.RewardJEXP),
		// TODO:
		// RewardItemList: dao.ToNullString(q.RewardItemList),
	}
}
