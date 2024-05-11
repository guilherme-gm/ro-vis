import { ArrayEqual } from "../../CompareUtils/ArrayEqual.js";
import { RecordObject } from "../../Database/RecordObject.js";
import { QuestRewardItem } from "./QuestRewardItem.js";

/**
 * Represents a Quest in the tool.
 */
export class Quest implements RecordObject {
	/**
	 * The File Version that originated this object
	 */
	public _FileVersion: number = -1;

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
	 * Always "QUE_NOIMAGE" or empty string (""), used to be something when viewing the quest,
	 * but was not really used officially
	 */
	public OldImage: string = "";

	/**
	 * Icon shown in the quest list.
	 * In the old UI (before 2018-03-21), that was always SG_FEEL,
	 * In the new UI, this uses a few images that represents "type of quest"
	 */
	public IconName: string = "";

	/**
	 * NPC Sprite
	 */
	public NpcSpr: string = "";

	/**
	 * NPC Map
	 */
	public NpcNavi: string = "";

	/**
	 * NPC X
	 */
	public NpcPosX: number = -1;

	/**
	 * NPC Y
	 */
	public NpcPosY: number = -1;

	/**
	 * Rewarded Base EXP
	 */
	public RewardEXP: string = "";

	/**
	 * Rewarded Job EXP
	 */
	public RewardJEXP: string = "";

	/**
	 * Rewarded items
	 */
	public RewardItemList: QuestRewardItem[] = [];

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public equals(other: Quest): boolean {
		// _FileVersion is not checked, if 2 versions exists but the record is the same, we don't care.
		return (
			other.Title === this.Title
			&& ArrayEqual.isEqual(other.Description, this.Description)
			&& other.Summary === this.Summary
			&& other.IconName === this.IconName
			&& other.NpcSpr === this.NpcSpr
			&& other.NpcNavi === this.NpcNavi
			&& other.NpcPosX === this.NpcPosX
			&& other.NpcPosY === this.NpcPosY
			&& other.RewardEXP === this.RewardEXP
			&& other.RewardJEXP === this.RewardJEXP
			&& ArrayEqual.isEqual(this.RewardItemList, other.RewardItemList, (a, b) => a.equals(b))
		);
	}
}
