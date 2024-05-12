import { Expose, Type } from "class-transformer";
import { IFileEntry } from "../../CommonLoader/IFileEntry.js";
import { Quest } from "./Quest.js";
import { QuestRewardItem } from "./QuestRewardItem.js";
import { QuestV } from "./QuestV.js";
import { QuestV3RewardItem } from "./QuestV3RewardItem.js";
import { ConvertClass } from "../../Utils/ConvertClass.js";

/**
 * System/OngoingQuestInfoList_True.lub
 * Since 2018-03-21
 */
export class QuestV3 implements IFileEntry<Quest> {
	public static isV3(quest: QuestV): quest is QuestV3 {
		return quest._FileVersion === 3;
	}

	@Expose()
	public readonly _FileVersion: number = 3;

	/**
	 * Quest ID
	 */
	@Expose()
	public Id: number = 0;

	/**
	 * Quest title ("yellow" title in side UI, or entry name in quest UI list)
	 */
	@Expose()
	public Title: string = "";

	/**
	 * Quest long description inside quest UI list
	 */
	@Expose()
	public Description: string[] = [];

	/**
	 * Short, one-line, mission description. Shown in side UI.
	 */
	@Expose()
	public Summary: string = "";

	/**
	 * First seem on 2018-04-04
	 */
	@Expose()
	public IconName: string = "";

	/**
	 * First seem on 2018-04-18
	 */
	@Expose()
	public NpcSpr: string = "";
	@Expose()
	public NpcNavi: string = "";
	@Expose()
	public NpcPosX: number = -1;
	@Expose()
	public NpcPosY: number = -1;
	@Expose()
	public RewardEXP: string = "";
	@Expose()
	public RewardJEXP: string = "";
	@Expose()
	@Type(()=> QuestV3RewardItem)
	public RewardItemList: QuestV3RewardItem[] = [];

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public toEntity(): Quest {
		return ConvertClass.convert(this, Quest);
	}
}
