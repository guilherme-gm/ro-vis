import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from '../../../CommonParser/TokenTextTableParser.js';
import { Result } from '../../../CommonParser/Result.js';
import { Logger } from '../../../Logger.js';

type ReadFn<T> = keyof {
	[P in keyof TokenTextTableParser as ReturnType<TokenTextTableParser[P]> extends Result<T>? P : never]: unknown;
};

export abstract class ItemKeyValueTableParser<U> extends TokenTextTableParser {
	protected static async internalFromFile<T, U>(filePath: string, readFn: ReadFn<U>, ctor: new (content: string, readFn: ReadFn<U>) => T): Promise<T> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new ctor(content, readFn);
	}

	private readFunction: ReadFn<U>;

	protected constructor(content: string, readFunction: ReadFn<U>) {
		super(content);
		this.readFunction = readFunction;
	}

	public parse(): Promise<Map<number, U>> {
		const table = new Map<number, U>();

		while (!this.isEndOfFile()) {
			try {
				const id = this.readIntCell().unwrap();
				const value = this[this.readFunction]().unwrap() as U;

				table.set(id, value);
			} catch (error) {
				Logger.error('Failed while reading entry; skipping...', error);
			}
		}

		return Promise.resolve(table);
	}
}
