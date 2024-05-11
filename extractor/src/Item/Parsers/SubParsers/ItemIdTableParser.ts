import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from '../../../CommonParser/TokenTextTableParser.js';
import { Logger } from '../../../Logger.js';

export abstract class ItemIdTableParser extends TokenTextTableParser {
	protected static async internalFromFile<T>(filePath: string, ctor: new (content: string) => T): Promise<T> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new ctor(content);
	}

	public parse(): Promise<number[]> {
		const table: number[] = [];

		while (!this.isEndOfFile()) {
			try {
				const id = this.readIntCell().unwrap();

				// Official files broke the format at some places, this is a workaround
				// when a record is complete but there are a few trailling characters
				if (!this.isLineStart()) {
					this.consumeRestOfLine();
				}

				table.push(id);
			} catch (error) {
				Logger.error('Failed while reading entry; skipping...', error);
			}
		}

		return Promise.resolve(table);
	}
}
