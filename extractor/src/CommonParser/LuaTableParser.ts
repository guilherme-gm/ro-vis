import * as childProcess from "child_process";
import { Logger } from "../Logger.js";

export class LuaTableParser<T> {
	protected filePath: string;

	constructor(filePath: string) {
		this.filePath = filePath;
	}

	protected async extractLuaTable(tableName: string, forceKeyToTable: boolean, extraFiles: string[] = []): Promise<T> {
		const stdoutBuffers: any = [];
		const readPromise = new Promise<string>((resolve, reject) => {
			const process = childProcess.spawn(
				'lua/lua54.exe',
				[
					'lua/lua2json.lua',
					tableName,
					this.filePath.replace(/lub$/i, "lua"),
					forceKeyToTable ? "true" : "false",
					...extraFiles
				],
			);
			process.stdout.on('data', (chunk) => {
				stdoutBuffers.push(chunk);
			});

			process.stderr.on('data', (chunk) => {
				Logger.error(chunk?.toString() ?? chunk);
			})

			process.on('close', (code) => {
				if (code === 0) {
					resolve(Buffer.concat(stdoutBuffers).toString());
				} else {
					reject(new Error(`child process exited with code ${code}`));
				}
			});
		});

		const result = await readPromise;

		return JSON.parse(result);
	}

	protected fixArrayObjects(data: any): any {
		if (
			!Array.isArray(data)
			&& typeof(data) === 'object'
			&& Object.keys(data).length === 0
		) {
			return [];
		} else {
			return data ?? [];
		}
	}
}
