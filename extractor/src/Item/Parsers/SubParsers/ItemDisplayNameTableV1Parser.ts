import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from '../../../CommonParser/TokenTextTableParser.js';

export class ItemDisplayNameTableV1Parser extends TokenTextTableParser {
	public static async fromFile(filePath: string): Promise<ItemDisplayNameTableV1Parser> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new ItemDisplayNameTableV1Parser(content);
	}

	public parse(): Promise<Map<number, string>> {
		const displayNameTable = new Map<number, string>();

		while (!this.isEndOfFile()) {
			const id = this.readIntCell().unwrap();
			const name = this.readCell().unwrap();

			// Official files broke the format at some places, this is a workaround
			// when a record is complete but there are extra chars after name.
			// -------------
			// 123#Name#blabla
			//
			// 1234#
			// ------------
			// "blabla" is ignored so we can properly parse the number afterwards.
			if (!this.isLineStart()) {
				this.consumeRestOfLine();
			}

			displayNameTable.set(id, name);
		}

		return Promise.resolve(displayNameTable);
	}
}
