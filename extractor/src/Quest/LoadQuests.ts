import { BasicLoader } from "../CommonLoader/BasicLoader.js";
import { IDataLoader } from "../CommonLoader/IDataLoader.js";
import { IParser } from "../CommonParser/IParser.js";
import { Logger } from "../Logger.js";
import { Quest } from "./DataStructures/Quest.js";
import { QuestV } from "./DataStructures/QuestV.js";
import { QuestV1Parser } from "./Parsers/QuestV1Parser.js";
import { QuestV3Parser } from "./Parsers/QuestV3Parser.js";
import { QuestV4Parser } from "./Parsers/QuestV4Parser.js";
import { QuestDb } from "./QuestDb.js";
import { Update } from "../Updates/Update.js";

export class LoadQuests extends BasicLoader<Quest, QuestV> implements IDataLoader {
	public name: string = LoadQuests.name;

	constructor() {
		super(new QuestDb());
	}

	private getQuestDataVersion(patch: Update): number {
		const date = patch._id.substring(0, 10);
		if (date.localeCompare('2012-03-14') < 0) {
			/**
			 * ??? introduced the quest UI using questid2display.txt
			 */
			return 1;
		} else if (date.localeCompare('2018-03-21') < 0) {
			/**
			 * 2012-03-14 was introduced the use of lua files/quest, which contained a few extra info.
			 * I am ignoring these extra info for now, but will keep the diff here for documentation purposes.
			 */
			return 1; // Actually 2
		} else if (date.localeCompare('2020-08-05') < 0) {
			/**
			 * 2018-03-21 was introduced "System/OnGoingQuestInfoList_True.lub" and "System/RecommendedQuestInfoList_True.lub".
			 * This marks the end of "questid2display.txt" and lua files/quest/*.lua
			 */
			return 3;
		} else if (date.localeCompare('9999-01-01') < 0) {
			/**
			 * 2020-08-05 introduces "CoolTimeQuest" into "System/OnGoingQuestInfoList_True.lub".
			 * This now becomes "V4"
			 */
			return 4;
		}

		throw new Error(`Quest version for patch "${patch._id}" is not mapped.`);
	}

	public hasFileOfInterest(update: Update): boolean {
		const version = this.getQuestDataVersion(update);
		let fileName: string;
		if (version === 1) {
			fileName = "data/questid2display.txt";
		} else if (version === 3 || version == 4) {
			fileName = "system/ongoingquestinfolist_true.lub";
		} else {
			throw new Error(`Unsupported quest version "${version}"`);
		}

		const entry = update.updates.find((f) => f.file.toLocaleLowerCase().includes(fileName));
		if (!entry) {
			return false;
		}

		return true;
	}

	protected async getParser(update: Update): Promise<IParser<QuestV>> {
		const version = this.getQuestDataVersion(update);
		Logger.info(`Version: ${version}`);
		if (version === 1) {
			return QuestV1Parser.fromFile(this.getPath(update, 'data/questid2display.txt'));
		} else if (version === 3) {
			return QuestV3Parser.fromFile(this.getPath(update, 'System/OnGoingQuestInfoList_True.lub'));
		} else if (version === 4) {
			return QuestV4Parser.fromFile(this.getPath(update, 'System/OnGoingQuestInfoList_True.lub'));
		} else {
			throw new Error(`Unsupported quest version "${version}"`);
		}
	}
}
