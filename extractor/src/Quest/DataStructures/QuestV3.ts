import { ArrayEqual } from "../../CompareUtils/ArrayEqual.js";
import { RecordObject } from "../../Database/RecordObject.js";
import { Quest } from "./Quest.js";

/**
 * System/OngoingQuestInfoList_True.lub
 * Since 2018-03-21
 */
export class QuestV3 {
	public static isV3(quest: Quest): quest is QuestV3 {
		return quest._FileVersion === 3;
	}

	public readonly _FileVersion: number = 3;

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

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public hasChange(other: RecordObject): boolean {
		if (other.getFileVersion() !== this.getFileVersion()) {
			return true;
		}

		if (!(other instanceof QuestV3)) {
			throw new Error('Invalid type');
		}

		return (
			other.Id != this.Id
			|| other.Title != this.Title
			|| ArrayEqual.isEqual(other.Description, this.Description)
			|| other.Summary != this.Summary
			|| other.IconName != this.IconName
		);
	}
}
