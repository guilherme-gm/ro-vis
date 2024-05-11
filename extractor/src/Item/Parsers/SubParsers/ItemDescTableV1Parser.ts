import { ItemKeyValueTableParser } from './ItemKeyValueTableParser.js';

export class ItemDescTableV1Parser extends ItemKeyValueTableParser<string> {
	public static async fromFile(filePath: string): Promise<ItemDescTableV1Parser> {
		// @ts-ignore
		return ItemDescTableV1Parser.internalFromFile(filePath, "readMultilineCell", ItemDescTableV1Parser);
	}
}
