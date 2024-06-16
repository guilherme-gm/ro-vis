import path from "path";
import fs from "fs/promises";
import { LogRecord } from "../Database/LogRecord.js";
import { RecordObject } from "../Database/RecordObject.js";
import { SqlConverter } from "./SqlConverter.js";
import { Logger } from "../Logger.js";

export class LogRecordSqlConverter<T extends RecordObject> extends SqlConverter {

	public override async convert(table: string, data: LogRecord<T>[]): Promise<void> {
		const replaces: string[] = [];
		const historyReplaces: string[] = [];

		data.forEach((row) => {
			if (row.current.value === null) {
				// @TODO:
				Logger.error('Unhandled null current element');
			} else {
				const convertedPairs: [string, string][] = [];
				convertedPairs.push(['Patch', `"${row.current.patch}"`]);
				this.convertInto(row.current.value, convertedPairs);

				const sql = this.generateReplace(table, convertedPairs);
				replaces.push(sql);
			}

			row.previous.forEach((previous, idx) => {
				if (previous.value === null) {
					// @TODO:
					Logger.error('Unhandled null previous element');
					return;
				}

				const convertedPairs: [string, string][] = [];
				convertedPairs.push(['HistoryId', `"${row._id}_${previous.patch}"`]);
				convertedPairs.push(['Patch', `"${previous.patch}"`]);

				if (idx > 0) {
					convertedPairs.push(['PreviousId', `"${row._id}_${row.previous[idx - 1]!.patch}"`]);
				} else {
					convertedPairs.push(['PreviousId', 'NULL']);
				}

				this.convertInto(previous.value, convertedPairs);

				const sql = this.generateReplace(`${table}_history`, convertedPairs);
				historyReplaces.push(sql);
			});
		});

		await fs.mkdir(path.resolve("out", "sql"), { recursive: true });
		await fs.writeFile(path.resolve("out", "sql", `${table}.sql`), replaces.join('\n'));
		await fs.writeFile(path.resolve("out", "sql", `${table}_history.sql`), historyReplaces.join('\n'));
	}
}
