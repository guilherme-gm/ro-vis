export type QuestRewardItem  = {
	ItemID: number;

	ItemNum: number;
}

export type Quest = {
	id: string;

	patch: string;

	_FileVersion: number;

	QuestId: number;

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
	id: number;

	patch: string;

	Title: string;
};
