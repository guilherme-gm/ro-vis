import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import { TokenTextTableParser } from "../../CommonParser/TokenTextTableParser.js";
import { QuestV0 } from "../DataStructures/QuestV0.js";

export class QuestV0Parser extends TokenTextTableParser {
	public static async fromFile(filePath: string): Promise<QuestV0Parser> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new QuestV0Parser(content);
	}

	public parse(): QuestV0[] {
		const quests: QuestV0[] = [];

		while (!this.isEndOfFile()) {
			const quest = new QuestV0();
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

		return quests;
	}
}
