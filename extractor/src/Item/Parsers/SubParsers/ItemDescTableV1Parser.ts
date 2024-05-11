import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from '../../../CommonParser/TokenTextTableParser.js';

export class ItemDescTableV1Parser extends TokenTextTableParser {
	public static async fromFile(filePath: string): Promise<ItemDescTableV1Parser> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new ItemDescTableV1Parser(content);
	}

	public parse(): Promise<Map<number, string[]>> {
		const descTable = new Map<number, string[]>();

		while (!this.isEndOfFile()) {
			try {
				const id = this.readIntCell().unwrap();
				const desc = this.readMultilineCell().unwrap();

				// Official files broke the format at some places, this is a workaround
				// when a record is complete but the summary has a few extra chars. e.g.:
				// -------------
				// Summary#blabla
				//
				// 1234#
				// ------------
				// "blabla" is ignored so we can properly parse the number afterwards.
				// if (!this.isLineStart()) {
				// 	this.consumeRestOfLine();
				// }

				descTable.set(id, desc.split('\n'));
			} catch (error) {
				console.log('----- Failed while reading entry; skipping.... -----');
				console.error(error)
				console.log('----------------------------------------------------');
				console.log('');
			}
		}

		return Promise.resolve(descTable);
	}
}
