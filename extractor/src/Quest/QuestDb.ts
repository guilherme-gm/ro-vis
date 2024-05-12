import { LogRecordDao } from "../Database/LogRecordDao.js";
import { Quest } from "./DataStructures/Quest.js";
import { QuestRewardItem } from "./DataStructures/QuestRewardItem.js";

export class QuestDb extends LogRecordDao<Quest> {
	constructor() {
		super('quests');
	}

	protected override toInstance(data: Quest | null): Quest | null {
		if (data === null) {
			return null;
		}

		const q = new Quest();

		q.Id = data.Id;
		q.Title = data.Title;
		q.Summary = data.Summary;
		q.Description = data.Description;
		q.IconName = data.IconName;
		q.NpcSpr = data.NpcSpr;
		q.NpcPosX = data.NpcPosX;
		q.RewardEXP = data.RewardEXP;
		q.RewardJEXP = data.RewardJEXP;
		q.RewardItemList = data.RewardItemList.map((v) => {
			const r = new QuestRewardItem();
			r.ItemID = v.ItemID;
			r.ItemNum = v.ItemNum;

			return r;
		});

		return q;
	}
}
