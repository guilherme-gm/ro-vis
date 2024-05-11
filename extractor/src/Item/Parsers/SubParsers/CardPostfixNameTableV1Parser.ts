import { ItemIdTableParser } from './ItemIdTableParser.js';

export class CardPostfixNameTableV1Parser extends ItemIdTableParser {
	public static async fromFile(filePath: string): Promise<CardPostfixNameTableV1Parser> {
		// @ts-ignore
		return CardPostfixNameTableV1Parser.internalFromFile(filePath, CardPostfixNameTableV1Parser);
	}
}
