import { RecordObject } from "../../Database/RecordObject.js";

// questid2display.txt ; Since a long time ago (2008?)
export class QuestV1 implements RecordObject {
	public readonly _FileVersion: number = 1;

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
	public Description: string = "";

	/**
	 * Short, one-line, mission description. Shown in side UI.
	 */
	public Summary: string = "";

	/**
	 * Always "SG_FEEL", used to be the icon in the list of quests,
	 * but was not really used officially
	 */
	public OldIcon: string = "";

	/**
	 * Always "QUE_NOIMAGE" or empty string (""), used to be something when viewing the quest,
	 * but was not really used officially
	 */
	public OldImage: string = "";

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

		if (!(other instanceof QuestV1)) {
			throw new Error('Invalid type');
		}

		return (
			other.Id != this.Id
			|| other.Title != this.Title
			|| other.Description != this.Description
			|| other.Summary != this.Summary
			|| other.OldIcon != this.OldIcon
			|| other.OldImage != this.OldImage
		);
	}
}
