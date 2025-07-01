export type QuestRewardItem  = {
	ItemID: number;

	ItemNum: number;
}

export type Quest = {
	HistoryID: string;

	FileVersion: number;

	QuestID: number;

	Title: string;

	Description: string;

	Summary: string;

	OldImage: string;

	IconName: string;

	NpcSpr: string;

	NpcNavi: string;

	NpcPosX: number;

	NpcPosY: number;

	RewardEXP: string;

	RewardJEXP: string;

	RewardItemList: QuestRewardItem[];

	CoolTimeQuest: number;
}

export type MinQuest = {
	QuestID: number;

	LastUpdate: string;

	Title: string;
};
