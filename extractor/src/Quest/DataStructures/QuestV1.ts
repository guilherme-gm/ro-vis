import { RecordObject } from "../../Database/RecordObject.js";

// questid2display.txt ; Since a long time ago (2008?)
export class QuestV1 implements RecordObject {
	public readonly _FileVersion: number = 1;

	public Id: number = 0;

	public Title: string = "";

	public Description: string = "";

	public Summary: string = "";

	public OldIcon: string = "";

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
