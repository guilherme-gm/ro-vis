import { Expose } from "class-transformer";
import { IFileEntry } from "../../CommonLoader/IFileEntry.js";
import { Quest } from "./Quest.js";
import { QuestV } from "./QuestV.js";

// questid2display.txt ; Since a long time ago (2008?)
export class QuestV1 implements IFileEntry<Quest> {
	public static isV1(quest: QuestV): quest is QuestV1 {
		return quest._FileVersion === 1;
	}

	@Expose()
	public readonly _FileVersion: number = 1;

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
	public Description: string = "";

	/**
	 * Short, one-line, mission description. Shown in side UI.
	 */
	@Expose()
	public Summary: string = "";

	/**
	 * Always "SG_FEEL", used to be the icon in the list of quests,
	 * but was not really used officially
	 */
	@Expose()
	public OldIcon: string = "";

	/**
	 * Always "QUE_NOIMAGE" or empty string (""), used to be something when viewing the quest,
	 * but was not really used officially
	 */
	@Expose()
	public OldImage: string = "";

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
		q.Description = [this.Description];
		q.Summary = this.Summary;
		q.IconName = this.OldIcon;
		q.OldImage = this.OldImage;

		return q;
	}
}
