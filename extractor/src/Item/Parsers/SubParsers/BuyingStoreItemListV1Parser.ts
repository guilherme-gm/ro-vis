import { ItemIdTableParser } from './ItemIdTableParser.js';

export class BuyingStoreItemListV1Parser extends ItemIdTableParser {
	public static async fromFile(filePath: string): Promise<BuyingStoreItemListV1Parser> {
		// @ts-ignore
		return BuyingStoreItemListV1Parser.internalFromFile(filePath, BuyingStoreItemListV1Parser);
	}
}
