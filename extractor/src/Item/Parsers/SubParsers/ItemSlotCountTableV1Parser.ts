import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from '../../../CommonParser/TokenTextTableParser.js';

export class ItemSlotCountTableV1Parser extends TokenTextTableParser {
	public static async fromFile(filePath: string): Promise<ItemSlotCountTableV1Parser> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new ItemSlotCountTableV1Parser(content);
	}

	public parse(): Promise<Map<number, number>> {
		const slotTable = new Map<number, number>();

		while (!this.isEndOfFile()) {
			try {
				const itemId = this.readIntCell().unwrap();
				const slots = this.readIntCell().unwrap();
				slotTable.set(itemId, slots);
			} catch (error) {
				console.log('----- Failed while reading entry; skipping.... -----');
				console.error(error)
				console.log('----------------------------------------------------');
				console.log('');
			}
		}

		return Promise.resolve(slotTable);
	}
}
