import { LogRecord } from "../Database/LogRecord.js";
import { QuestV0 } from "./DataStructures/QuestV0.js";
import { QuestV0Parser } from "./Parsers/QuestV0Parser.js";
import { QuestDb } from "./QuestDb.js";
import { IDataLoader } from "../IDataLoader.js";
import { PatchRecord } from "../Patches/PatchRecord.js";
import path from "path";
import { patchesRootDir } from "../Config/config.js";
import * as fs from "fs";

export class LoadQuests implements IDataLoader {
	public name: string = LoadQuests.name;

	private getQuestDataVersion(patch: PatchRecord): number {
		const date = patch._id.substring(0, 10);
		if (date.localeCompare('2012-03-14') < 0) {
			/**
			 * 2012-03-14 was introduced the use of lua files/quest, which contained a few extra info.
			 * I am ignoring these extra info for now, but will keep the diff here for documentation purposes.
			 */
			return 1;
		} else if (date.localeCompare('2018-03-21') < 0) {
			/**
			 * 2018-03-21 was introduced "System/OnGoingQuestInfoList_True.lub" and "System/RecommendedQuestInfoList_True.lub".
			 * This marks the end of "questid2display.txt" and lua files/quest/*.lua
			 */
			return 1; // Actually 2 due to previous if
		} /* else if (date.localeCompare('9999-12-31') < 0) {
			return 3;
		} */

		throw new Error(`Quest version for patch "${patch._id}" is not mapped.`);
	}

	public hasFileOfInterest(patch: PatchRecord): boolean {
		const version = this.getQuestDataVersion(patch);
		if (version > 1) {
			console.warn(`Unsupported quest version "${version}"`);
			return false;
		}

		const entry = patch.files.find((f) => f.toLocaleLowerCase().includes("data\\questid2display.txt"));
		if (!entry) {
			return false;
		}

		// "[N]o changes"
		if (entry.startsWith("[N]")) {
			return false;
		}

		return true;
	}

	private async getQuestList(patchFolder: string): Promise<QuestV0[]> {
		const parser = await QuestV0Parser.fromFile(path.join(patchFolder, 'data', 'questid2display.txt'));
		const rawQuests = parser.parse();

		const questMap = new Map<string, QuestV0>();
		rawQuests.forEach((quest) => {
			questMap.set(quest.getId(), quest);
		});

		return [...questMap.values()];
	}

	public async load(patch: PatchRecord): Promise<void> {
		const patchFolder = path.join(patchesRootDir, patch._id);
		if (!fs.existsSync(patchFolder)) {
			console.warn(`!!!! WARN: Folder not found patch "${patch._id}" for file "questid2display"`);
			fs.appendFileSync("./not-found.txt", `${patch._id}\tdata/questid2display.txt\n`);
			return;
		}

		const quests = await this.getQuestList(patchFolder);

		const questDb = new QuestDb();
		const existingRecords = (await questDb.getAll()).reduce(
			(memo, record) => {
				memo.set(record._id, record);
				return memo;
			},
			new Map<string, LogRecord<QuestV0>>()
		);

		const newRecords: Map<string, LogRecord<QuestV0>> = new Map<string, LogRecord<QuestV0>>();
		const updatedRecords: LogRecord<QuestV0>[] = [];

		for (const quest of quests) {
			const record = existingRecords.get(quest.getId());
			if (!record) {
				newRecords.set(quest.getId(), new LogRecord<QuestV0>(patch._id, quest));
			} else {
				if (record.current.value.hasChange(quest)) {
					record.addChange(patch._id, quest);
					updatedRecords.push(record);
				}
			}
		}

		// fs.writeFileSync(`out_${patch._id}_new.json`, JSON.stringify([...newRecords.values()], null, 4));
		// fs.writeFileSync(`out_${patch._id}_upd.json`, JSON.stringify([...updatedRecords], null, 4));

		if (newRecords.size === 0 && updatedRecords.length === 0) {
			console.warn(`!!!! WARN: NO-Change patch "${patch._id}" for file "questid2display"`);
			fs.appendFileSync("./no-op.txt", `${patch._id}\tdata/questid2display.txt\n`);
			return;
		}

		console.log(`${newRecords.size} new records to create and ${updatedRecords.length} to update...`);
		const newRecordsArr = [...newRecords.values()];
		while (newRecordsArr.length > 0) {
			await questDb.insertMany(newRecordsArr.splice(0, 500));
		}

		for (let i = 0; i < updatedRecords.length; i++) {
			await questDb.updateOrCreate(updatedRecords[i]!);
		}
	}
}
