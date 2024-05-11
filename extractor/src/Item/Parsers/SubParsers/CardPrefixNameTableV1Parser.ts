import { ItemKeyValueTableParser } from './ItemKeyValueTableParser.js';

export class CardPrefixNameTableV1Parser extends ItemKeyValueTableParser<string> {
	public static async fromFile(filePath: string): Promise<CardPrefixNameTableV1Parser> {
		// @ts-ignore
		return CardPrefixNameTableV1Parser.internalFromFile(filePath, "readCell", CardPrefixNameTableV1Parser);
	}
}
