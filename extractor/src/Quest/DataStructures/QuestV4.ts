import { IFileEntry } from "../../CommonLoader/IFileEntry.js";
import { Quest } from "./Quest.js";
import { QuestRewardItem } from "./QuestRewardItem.js";
import { QuestV } from "./QuestV.js";
import { QuestV3RewardItem } from "./QuestV3RewardItem.js";

/**
 * System/OngoingQuestInfoList_True.lub
 * Since 2020-08-05
 * V4 - Adds CoolTimeQuest
 */
export class QuestV4 implements IFileEntry<Quest> {
	public static isV4(quest: QuestV): quest is QuestV4 {
		return quest._FileVersion === 4;
	}

	public readonly _FileVersion: number = 4;

	/**
	 * Quest ID
	 */
	public Id: number = 0;

	/**
	 * Quest title ("yellow" title in side UI, or entry name in quest UI list)
	 */
	public Title: string = "";

	/**
	 * Quest long description inside quest UI list
	 */
	public Description: string[] = [];

	/**
	 * Short, one-line, mission description. Shown in side UI.
	 */
	public Summary: string = "";

	/**
	 * First seem on 2018-04-04
	 */
	public IconName: string = "";

	/**
	 * First seem on 2018-04-18
	 */
	public NpcSpr: string = "";
	public NpcNavi: string = "";
	public NpcPosX: number = -1;
	public NpcPosY: number = -1;
	public RewardEXP: string = "";
	public RewardJEXP: string = "";
	public RewardItemList: QuestV3RewardItem[] = [];

	/**
	 * Is it a Cooldown quest (?)
	 * Since 2020-08-05
	 */
	public CoolTimeQuest: number = 0;

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public toEntity(): Quest {
		const q = new Quest();

		q._FileVersion = this._FileVersion;
		q.Id = this.Id;
		q.Title = this.Title;
		q.Description = this.Description;
		q.Summary = this.Summary;
		q.IconName = this.IconName;
		q.NpcSpr = this.NpcSpr;
		q.NpcNavi = this.NpcNavi;
		q.NpcPosX = this.NpcPosX;
		q.NpcPosY = this.NpcPosY;
		q.RewardEXP = this.RewardEXP;
		q.RewardJEXP = this.RewardJEXP;
		q.CoolTimeQuest = this.CoolTimeQuest;

		q.RewardItemList = [];
		this.RewardItemList?.forEach((r) => {
			const reward = new QuestRewardItem();
			reward.ItemID = r.ItemID;
			reward.ItemNum = r.ItemNum;

			q.RewardItemList.push(reward);
		});

		return q;
	}
}
