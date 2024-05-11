import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from '../../../CommonParser/TokenTextTableParser.js';
import { ItemKeyValueTableParser } from './ItemKeyValueTableParser.js';

export class ItemSlotCountTableV1Parser extends ItemKeyValueTableParser<number> {
	public static async fromFile(filePath: string): Promise<ItemSlotCountTableV1Parser> {
		// @ts-ignore
		return ItemSlotCountTableV1Parser.internalFromFile(filePath, "readIntCell", ItemSlotCountTableV1Parser);
	}
}
