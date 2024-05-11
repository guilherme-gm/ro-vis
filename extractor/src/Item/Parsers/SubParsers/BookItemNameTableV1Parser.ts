import { ItemIdTableParser } from './ItemIdTableParser.js';

export class BookItemNameTableV1Parser extends ItemIdTableParser {
	public static async fromFile(filePath: string): Promise<BookItemNameTableV1Parser> {
		// @ts-ignore
		return BookItemNameTableV1Parser.internalFromFile(filePath, BookItemNameTableV1Parser);
	}
}
