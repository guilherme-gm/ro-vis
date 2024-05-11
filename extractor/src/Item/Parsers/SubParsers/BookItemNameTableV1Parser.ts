import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from '../../../CommonParser/TokenTextTableParser.js';

export class BookItemNameTableV1Parser extends TokenTextTableParser {
	public static async fromFile(filePath: string): Promise<BookItemNameTableV1Parser> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new BookItemNameTableV1Parser(content);
	}

	public parse(): Promise<number[]> {
		const bookItems: number[] = [];

		while (!this.isEndOfFile()) {
			try {
				const id = this.readIntCell().unwrap();
				bookItems.push(id);
			} catch (error) {
				console.log('----- Failed while reading entry; skipping.... -----');
				console.error(error)
				console.log('----------------------------------------------------');
				console.log('');
			}
		}

		return Promise.resolve(bookItems);
	}
}
