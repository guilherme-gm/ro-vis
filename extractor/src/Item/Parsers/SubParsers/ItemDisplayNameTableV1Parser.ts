import { ItemKeyValueTableParser } from './ItemKeyValueTableParser.js';

export class ItemDisplayNameTableV1Parser extends ItemKeyValueTableParser<string> {
	public static async fromFile(filePath: string): Promise<ItemDisplayNameTableV1Parser> {
		// @ts-ignore
		return ItemDisplayNameTableV1Parser.internalFromFile(filePath, "readCell", ItemDisplayNameTableV1Parser);
	}
}
