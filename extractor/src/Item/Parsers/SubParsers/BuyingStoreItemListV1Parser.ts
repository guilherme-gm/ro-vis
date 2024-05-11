import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from '../../../CommonParser/TokenTextTableParser.js';

export class BuyingStoreItemListV1Parser extends TokenTextTableParser {
	public static async fromFile(filePath: string): Promise<BuyingStoreItemListV1Parser> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new BuyingStoreItemListV1Parser(content);
	}

	public parse(): Promise<number[]> {
		const buyingStoreItems: number[] = [];

		while (!this.isEndOfFile()) {
			try {
				const id = this.readIntCell().unwrap();
				buyingStoreItems.push(id);
			} catch (error) {
				console.log('----- Failed while reading entry; skipping.... -----');
				console.error(error)
				console.log('----------------------------------------------------');
				console.log('');
			}
		}

		return Promise.resolve(buyingStoreItems);
	}
}
