import * as fs from "fs";
import chalk from "chalk";
import path from "path";
import { IParser } from "../CommonParser/IParser.js";
import { LogRecord } from "../Database/LogRecord.js";
import { IDataLoader } from "../IDataLoader.js";
import { PatchRecord } from "../Patches/PatchRecord.js";
import { ItemV } from "./DataStructures/ItemV.js";
import { Item } from "./DataStructures/Item.js";
import { ItemDb } from "./ItemDb.js";
import { ItemV1Parser } from "./Parsers/ItemV1Parser.js";
import { Logger } from "../Logger.js";
import { Config } from "../Config/config.js";
import { Cli } from "../Cli.js";

export class LoadItem implements IDataLoader {
	public name: string = LoadItem.name;

	private readonly v1FileNames = [
		"data\\bookitemnametable.txt",
		"data\\buyingstoreitemlist.txt",
		"data\\cardpostfixnametable.txt",
		"data\\cardprefixnametable.txt",
		"data\\idnum2itemdesctable.txt",
		"data\\idnum2itemdisplaynametable.txt",
		"data\\idnum2itemresnametable.txt",
		"data\\itemslotcounttable.txt",
		"data\\num2cardillustnametable.txt",
		"data\\num2itemdesctable.txt",
		"data\\num2itemdisplaynametable.txt",
		"data\\num2itemresnametable.txt",
	];

	private getItemDataVersion(patch: PatchRecord): number {
		const date = patch._id.substring(0, 10);
		if (date.localeCompare('2012-04-18') < 0) {
			/**
			 * Everything is txt
			 */
			return 1;
		}

		throw new Error(`Item version for patch "${patch._id}" is not mapped.`);
	}

	public hasFileOfInterest(patch: PatchRecord): boolean {
		const version = this.getItemDataVersion(patch);
		let fileNames: string[];
		if (version === 1) {
			fileNames = this.v1FileNames;
		} else {
			throw new Error(`Unsupported item version "${version}"`);
		}

		fileNames = fileNames.map((n) => n.toLocaleLowerCase());

		const entry = patch.files.find((f) => fileNames.includes(f.toLocaleLowerCase()));
		if (!entry) {
			return false;
		}

		return true;
	}

	private getPathIfExists(patch: PatchRecord, patchFolder: string, file: string): string | null {
		const desiredFile = patch.files.find((f) => f.toLocaleLowerCase().includes(file.toLocaleLowerCase()));
		if (!desiredFile) {
			return null;
		}

		return path.join(patchFolder, desiredFile);
	}

	private async getParser(patch: PatchRecord, patchFolder: string, itemMap: Map<number, Item>): Promise<IParser<ItemV>> {
		const version = this.getItemDataVersion(patch);
		Logger.info(`${chalk.whiteBright('Version:')} ${version}`);
		if (version === 1) {
			return new ItemV1Parser(itemMap, {
				bookItemNameTable: this.getPathIfExists(patch, patchFolder, "data\\bookitemnametable.txt"),
				buyingStoreItemList: this.getPathIfExists(patch, patchFolder, "data\\buyingstoreitemlist.txt"),
				cardPostFixNameTable: this.getPathIfExists(patch, patchFolder, "data\\cardpostfixnametable.txt"),
				cardPrefixNameTable: this.getPathIfExists(patch, patchFolder, "data\\cardprefixnametable.txt"),
				idNum2ItemDescTable: this.getPathIfExists(patch, patchFolder, "data\\idnum2itemdesctable.txt"),
				idNum2ItemDisplayNameTable: this.getPathIfExists(patch, patchFolder, "data\\idnum2itemdisplaynametable.txt"),
				idNum2ItemResNameTable: this.getPathIfExists(patch, patchFolder, "data\\idnum2itemresnametable.txt"),
				itemSlotCountTable: this.getPathIfExists(patch, patchFolder, "data\\itemslotcounttable.txt"),
				num2CardIllustNameTable: this.getPathIfExists(patch, patchFolder, "data\\num2cardillustnametable.txt"),
				num2ItemDescTable: this.getPathIfExists(patch, patchFolder, "data\\num2itemdesctable.txt"),
				num2ItemDisplayNameTable: this.getPathIfExists(patch, patchFolder, "data\\num2itemdisplaynametable.txt"),
				num2ItemResNameTable: this.getPathIfExists(patch, patchFolder, "data\\num2itemresnametable.txt"),
			});
		} else {
			throw new Error(`Unsupported quest version "${version}"`);
		}
	}

	private async getItemList(patch: PatchRecord, patchFolder: string, itemMap: Map<number, Item>): Promise<Item[]> {
		const parser = await this.getParser(patch, patchFolder, itemMap);
		const rawItems = await parser.parse();

		const itemsMap = new Map<string, Item>();
		rawItems.forEach((item) => {
			itemsMap.set(item.getId(), item.toItem());
		});

		return [...itemsMap.values()];
	}

	public async load(patch: PatchRecord): Promise<void> {
		const itemDb = new ItemDb();
		if (Cli.cli.dryRun && !Cli.cli.cleanRun) {
			await itemDb.replicate();
		}

		const existingRecords = (await itemDb.getAll()).reduce(
			(memo, record) => {
				memo.set(record._id, record);
				return memo;
			},
			new Map<string, LogRecord<Item>>()
		);

		const itemMap = new Map<number, Item>();
		existingRecords.forEach((item) => {
			itemMap.set(item.current.value.Id, item.current.value);
		});

		const patchFolder = path.join(Config.patchesRootDir, patch._id);
		if (!fs.existsSync(patchFolder)) {
			Logger.warn(`!!!! WARN: Folder not found patch "${patch._id}" for file "questid2display"`);
			fs.appendFileSync("./not-found.txt", `${patch._id}\tdata/questid2display.txt\n`);
			return;
		}

		const items = await this.getItemList(patch, patchFolder, itemMap);
		const newRecords: Map<string, LogRecord<Item>> = new Map<string, LogRecord<Item>>();
		const updatedRecords: LogRecord<Item>[] = [];

		for (const item of items) {
			const record = existingRecords.get(item.getId());
			if (!record) {
				newRecords.set(item.getId(), new LogRecord<Item>(patch._id, item));
			} else {
				if (!record.current.value.equals(item)) {
					record.addChange(patch._id, item);
					updatedRecords.push(record);
				}
			}
		}

		// fs.writeFileSync(`out_${patch._id}_new.json`, JSON.stringify([...newRecords.values()], null, 4));
		// fs.writeFileSync(`out_${patch._id}_upd.json`, JSON.stringify([...updatedRecords], null, 4));

		if (newRecords.size === 0 && updatedRecords.length === 0) {
			Logger.warn(`!!!! WARN: NO-Change patch "${patch._id}" for file "questid2display"`);
			fs.appendFileSync("./no-op.txt", `${patch._id}\tdata/questid2display.txt\n`);
			return;
		}

		Logger.info(`${newRecords.size} new records to create and ${updatedRecords.length} to update...`);
		const newRecordsArr = [...newRecords.values()];
		while (newRecordsArr.length > 0) {
			await itemDb.insertMany(newRecordsArr.splice(0, 500));
		}

		if (updatedRecords.length > 0) {
			await itemDb.bulkWrite(updatedRecords);
		}
	}
}
