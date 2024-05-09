import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from "../../CommonParser/TokenTextTableParser.js";
import { QuestV1 } from "../DataStructures/QuestV1.js";

export class QuestV1Parser extends TokenTextTableParser {
	public static async fromFile(filePath: string): Promise<QuestV1Parser> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new QuestV1Parser(content);
	}

	public parse(): Promise<QuestV1[]> {
		const quests: QuestV1[] = [];

		while (!this.isEndOfFile()) {
			const quest = new QuestV1();
			quest.Id = this.readIntCell().unwrap();
			quest.Title = this.readCell().unwrap();
			quest.OldIcon = this.readCell().unwrap();
			quest.OldImage = this.readCell().unwrap();
			quest.Description = this.readMultilineCell().unwrap();
			quest.Summary = this.readMultilineCell().unwrap();

			// Official files broke the format at some places, this is a workaround
			// when a record is complete but the summary has a few extra chars. e.g.:
			// -------------
			// Summary#blabla
			//
			// 1234#
			// ------------
			// "blabla" is ignored so we can properly parse the number afterwards.
			if (!this.isLineStart()) {
				this.consumeRestOfLine();
			}

			quests.push(quest);
		}

		return Promise.resolve(quests);
	}
}
