import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from '../../../CommonParser/TokenTextTableParser.js';

export class ItemResNameTableV1Parser extends TokenTextTableParser {
	public static async fromFile(filePath: string): Promise<ItemResNameTableV1Parser> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new ItemResNameTableV1Parser(content);
	}

	public parse(): Promise<Map<number, string>> {
		const resNameTable = new Map<number, string>();

		while (!this.isEndOfFile()) {
			const id = this.readIntCell().unwrap();
			const name = this.readCell().unwrap();

			// Official files broke the format at some places, this is a workaround
			// when a record is complete but there are extra chars in the line
			if (!this.isLineStart()) {
				this.consumeRestOfLine();
			}

			resNameTable.set(id, name);
		}

		return Promise.resolve(resNameTable);
	}
}
