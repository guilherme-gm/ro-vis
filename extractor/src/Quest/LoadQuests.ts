import * as fs from "fs";
import path from "path";
import { IParser } from "../CommonParser/IParser.js";
import { LogRecord } from "../Database/LogRecord.js";
import { IDataLoader } from "../IDataLoader.js";
import { PatchRecord } from "../Patches/PatchRecord.js";
import { Quest } from "./DataStructures/Quest.js";
import { QuestV } from "./DataStructures/QuestV.js";
import { QuestV1Parser } from "./Parsers/QuestV1Parser.js";
import { QuestV3Parser } from "./Parsers/QuestV3Parser.js";
import { QuestDb } from "./QuestDb.js";
import { QuestV4Parser } from "./Parsers/QuestV4Parser.js";
import { Logger } from "../Logger.js";
import { Config } from "../Config/config.js";
import { Cli } from "../Cli.js";

export class LoadQuests implements IDataLoader {
	public name: string = LoadQuests.name;

	private getQuestDataVersion(patch: PatchRecord): number {
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

	public hasFileOfInterest(patch: PatchRecord): boolean {
		const version = this.getQuestDataVersion(patch);
		let fileName: string;
		if (version === 1) {
			fileName = "data\\questid2display.txt";
		} else if (version === 3 || version == 4) {
			fileName = "system\\ongoingquestinfolist_true.lub";
		} else {
			throw new Error(`Unsupported quest version "${version}"`);
		}

		const entry = patch.files.find((f) => f.toLocaleLowerCase().includes(fileName));
		if (!entry) {
			return false;
		}

		// "[N]o changes"
		if (entry.startsWith("[N]")) {
			return false;
		}

		return true;
	}

	private async getParser(patch: PatchRecord, patchFolder: string): Promise<IParser<QuestV>> {
		const version = this.getQuestDataVersion(patch);
		Logger.info(`Version: ${version}`);
		if (version === 1) {
			return QuestV1Parser.fromFile(path.join(patchFolder, 'data', 'questid2display.txt'));
		} else if (version === 3) {
			return QuestV3Parser.fromFile(path.join(patchFolder, 'System', 'OnGoingQuestInfoList_True.lub'));
		} else if (version === 4) {
			return QuestV4Parser.fromFile(path.join(patchFolder, 'System', 'OnGoingQuestInfoList_True.lub'));
		} else {
			throw new Error(`Unsupported quest version "${version}"`);
		}
	}

	private async getQuestList(patch: PatchRecord, patchFolder: string): Promise<Quest[]> {
		const parser = await this.getParser(patch, patchFolder);
		const rawQuests = await parser.parse();

		const questMap = new Map<string, Quest>();
		rawQuests.forEach((quest) => {
			questMap.set(quest.getId(), quest.toQuest());
		});

		return [...questMap.values()];
	}

	public async load(patch: PatchRecord): Promise<void> {
		const patchFolder = path.join(Config.patchesRootDir, patch._id);
		if (!fs.existsSync(patchFolder)) {
			Logger.warn(`!!!! WARN: Folder not found patch "${patch._id}" for file "questid2display"`);
			fs.appendFileSync("./not-found.txt", `${patch._id}\tdata/questid2display.txt\n`);
			return;
		}

		const quests = await this.getQuestList(patch, patchFolder);

		const questDb = new QuestDb();
		if (Cli.cli.dryRun && !Cli.cli.cleanRun) {
			await questDb.replicate();
		}

		const existingRecords = (await questDb.getAll()).reduce(
			(memo, record) => {
				memo.set(record._id, record);
				return memo;
			},
			new Map<string, LogRecord<Quest>>()
		);

		const newRecords: Map<string, LogRecord<Quest>> = new Map<string, LogRecord<Quest>>();
		const updatedRecords: LogRecord<Quest>[] = [];

		for (const quest of quests) {
			const record = existingRecords.get(quest.getId());
			if (!record) {
				newRecords.set(quest.getId(), new LogRecord<Quest>(patch._id, quest));
			} else {
				if (!record.current.value.equals(quest)) {
					record.addChange(patch._id, quest);
					updatedRecords.push(record);
				}
			}
		}

		if (Cli.cli.dryRun) {
			fs.writeFileSync(`out/out_${patch._id}_${this.name}_new.json`, JSON.stringify([...newRecords.values()], null, 4));
			fs.writeFileSync(`out/out_${patch._id}_${this.name}_upd.json`, JSON.stringify([...updatedRecords], null, 4));
		}

		if (newRecords.size === 0 && updatedRecords.length === 0) {
			Logger.warn(`!!!! WARN: NO-Change patch "${patch._id}" for file "questid2display"`);
			fs.appendFileSync("./no-op.txt", `${patch._id}\tdata/questid2display.txt\n`);
			return;
		}

		Logger.info(`${newRecords.size} new records to create and ${updatedRecords.length} to update...`);
		const newRecordsArr = [...newRecords.values()];
		while (newRecordsArr.length > 0) {
			await questDb.insertMany(newRecordsArr.splice(0, 500));
		}

		if (updatedRecords.length > 0) {
			await questDb.bulkWrite(updatedRecords);
		}
	}
}
