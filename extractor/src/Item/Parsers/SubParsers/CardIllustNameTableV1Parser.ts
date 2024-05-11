import { ItemKeyValueTableParser } from './ItemKeyValueTableParser.js';

export class CardIllustNameTableV1Parser extends ItemKeyValueTableParser<string> {
	public static async fromFile(filePath: string): Promise<CardIllustNameTableV1Parser> {
		// @ts-ignore
		return CardIllustNameTableV1Parser.internalFromFile(filePath, "readCell", CardIllustNameTableV1Parser);
	}
}
