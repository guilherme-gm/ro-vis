import { ItemKeyValueTableParser } from './ItemKeyValueTableParser.js';

export class ItemResNameTableV1Parser extends ItemKeyValueTableParser<string> {
	public static async fromFile(filePath: string): Promise<ItemResNameTableV1Parser> {
		// @ts-ignore
		return ItemResNameTableV1Parser.internalFromFile(filePath, "readCell", ItemResNameTableV1Parser);
	}
}
