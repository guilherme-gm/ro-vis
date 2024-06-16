export type QuestRewardItem  = {
	ItemID: number;

	ItemNum: number;
}

export type Quest = {
	HistoryId: string;

	Patch: string;

	_FileVersion: number;

	Id: number;

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
	Id: number;

	Patch: string;

	Title: string;
};
