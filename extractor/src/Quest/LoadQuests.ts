import { LogRecord } from "../Database/LogRecord.js";
import { QuestV0 } from "./DataStructures/QuestV0.js";
import { QuestV0Parser } from "./Parsers/QuestV0Parser.js";
import { QuestDb } from "./QuestDb.js";
import { IDataLoader } from "../IDataLoader.js";
import { PatchRecord } from "../Patches/PatchRecord.js";

export class LoadQuests implements IDataLoader {
	public hasFileOfInterest(patch: PatchRecord): boolean {
		return patch.files.some((f) => f.toLocaleLowerCase().includes("data\\questid2display.txt"));
	}

	public async load(patch: PatchRecord): Promise<void> {
		// @TODO

		// await questDb.updateOrCreate(newRecords);
	}
}
