import * as fs from "fs";
import { Cli } from "../Cli.js";
import { IDataLoader } from "../CommonLoader/IDataLoader.js";
import { IParser, ParserResult } from "../CommonParser/IParser.js";
import { LogRecord } from "../Database/LogRecord.js";
import { LogRecordDao } from "../Database/LogRecordDao.js";
import { RecordObject } from "../Database/RecordObject.js";
import { Logger } from "../Logger.js";
import { IFileEntry } from "./IFileEntry.js";
import { ParsingResult } from "../CommonParser/ParsingResult.js";
import { Update } from "../Updates/Update.js";
import { Config } from "../Config/config.js";
import path from "path";

export abstract class BasicLoader<T extends RecordObject, U extends IFileEntry<T>> implements IDataLoader {
	public abstract name: string;

	protected entityDb: LogRecordDao<T>;

	protected existingRecords: Map<string, LogRecord<T>> = new Map<string, LogRecord<T>>();

	protected constructor(entityDb: LogRecordDao<T>) {
		this.entityDb = entityDb;
	}

	protected getPathIfExists(update: Update, file: string): string | null {
		const desiredFile = update.updates.find((f) => f.file.toLocaleLowerCase().includes(file.toLocaleLowerCase()));
		if (!desiredFile) {
			return null;
		}

		return path.join(Config.patchesRootDir, desiredFile.patch, desiredFile.file);
	}

	protected getPath(update: Update, file: string): string {
		const filePath = this.getPathIfExists(update, file);
		if (!filePath) {
			throw new Error(`File ${file} does not exists in update ${update._id} file list.`);
		}

		if (!fs.existsSync(filePath)) {
			throw new Error(`File ${filePath} does not exists in disk.`);
		}

		return filePath;
	}

	public abstract hasFileOfInterest(update: Update): boolean;

	protected abstract getParser(update: Update): Promise<IParser<U>>;

	protected async getPatchEntries(update: Update): Promise<ParserResult<T>> {
		const parser = await this.getParser(update);
		const rawEntries = await parser.parse();

		if (rawEntries.result === ParsingResult.Empty) {
			return {
				result: ParsingResult.Empty,
				data: [],
			};
		}

		const entriesMap = new Map<string, T>();
		rawEntries.data.forEach((entry) => {
			entriesMap.set(entry.getId(), entry.toEntity());
		});

		return {
			result: ParsingResult.Ok,
			data: [...entriesMap.values()],
		};
	}

	public async load(update: Update): Promise<void> {
		if (Cli.cli.dryRun) {
			await this.entityDb.replicate();
		}

		this.existingRecords = (await this.entityDb.getAll()).reduce(
			(memo, record) => {
				memo.set(record._id, record);
				return memo;
			},
			new Map<string, LogRecord<T>>()
		);

		const patchEntries = await this.getPatchEntries(update);
		if (patchEntries.result === ParsingResult.Empty) {
			Logger.warn(`The patch doesn't have files to load. (Probably same file)`);
			return;
		}

		const newRecords: Map<string, LogRecord<T>> = new Map<string, LogRecord<T>>();
		const updatedRecords: LogRecord<T>[] = [];

		const patchIds = new Set<string>();
		for (const patchEntry of patchEntries.data) {
			patchIds.add(patchEntry.getId());
			const record = this.existingRecords.get(patchEntry.getId());
			if (!record) {
				newRecords.set(patchEntry.getId(), new LogRecord<T>(update._id, patchEntry));
			} else {
				const currentValue = record.current.value;
				if (currentValue === null || !currentValue.equals(patchEntry)) {
					record.addChange(update._id, patchEntry);
					updatedRecords.push(record);
				}
			}
		}

		for (const existingRecord of this.existingRecords.values()) {
			if (!patchIds.has(existingRecord._id)) {
				existingRecord.addChange(update._id, null);
			}
		}

		if (Cli.cli.changeDump) {
			fs.writeFileSync(`out/out_${update._id}_${this.name}_new.json`, JSON.stringify([...newRecords.values()], null, 4));
			fs.writeFileSync(`out/out_${update._id}_${this.name}_upd.json`, JSON.stringify([...updatedRecords], null, 4));
		}

		if (newRecords.size === 0 && updatedRecords.length === 0) {
			Logger.info('The patch does not cause record changes. (Probably same file)');
			return;
		}

		Logger.info(`${newRecords.size} new records to create and ${updatedRecords.length} to update...`);
		const newRecordsArr = [...newRecords.values()];
		while (newRecordsArr.length > 0) {
			await this.entityDb.insertMany(newRecordsArr.splice(0, 500));
		}

		if (updatedRecords.length > 0) {
			await this.entityDb.bulkWrite(updatedRecords);
		}
	}
}
