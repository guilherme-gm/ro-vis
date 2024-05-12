import { Expose } from "class-transformer";
import { IFileEntry } from "../../CommonLoader/IFileEntry.js";
import { Quest } from "./Quest.js";
import { QuestV } from "./QuestV.js";
import { QuestV3RewardItem } from "./QuestV3RewardItem.js";
import { ConvertClass } from "../../Utils/ConvertClass.js";

/**
 * System/OngoingQuestInfoList_True.lub
 * Since 2020-08-05
 * V4 - Adds CoolTimeQuest
 */
export class QuestV4 implements IFileEntry<Quest> {
	public static isV4(quest: QuestV): quest is QuestV4 {
		return quest._FileVersion === 4;
	}

	@Expose()
	public readonly _FileVersion: number = 4;

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
	public RewardItemList: QuestV3RewardItem[] = [];

	/**
	 * Is it a Cooldown quest (?)
	 * Since 2020-08-05
	 */
	@Expose()
	public CoolTimeQuest: number = 0;

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
