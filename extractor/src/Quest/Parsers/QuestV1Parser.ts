import iconv from 'iconv-lite';
import * as fs from 'fs/promises';
import * as fsSync from 'fs';
import { TokenTextTableParser } from "../../CommonParser/TokenTextTableParser.js";
import { QuestV1 } from "../DataStructures/QuestV1.js";
import { Logger } from '../../Logger.js';
import { ParsingResult } from '../../CommonParser/ParsingResult.js';
import { ParserResult } from '../../CommonParser/IParser.js';

export class QuestV1Parser extends TokenTextTableParser {
	public static async fromFile(filePath: string): Promise<QuestV1Parser> {
		if (!fsSync.existsSync(filePath)) {
			return new QuestV1Parser('', true);
		}

		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new QuestV1Parser(content, false);
	}

	private isEmpty: boolean = false;

	constructor(content: string, isEmpty: boolean) {
		super(content);
		this.isEmpty = isEmpty;
	}

	public parse(): Promise<ParserResult<QuestV1>> {
		if (this.isEmpty) {
			return Promise.resolve({
				result: ParsingResult.Empty,
				data: [],
			});
		}

		const quests: QuestV1[] = [];

		while (!this.isEndOfFile()) {
			try {
				const quest = new QuestV1();
				quest.Id = this.readIntCell().unwrap();
				quest.Title = this.readCell().unwrap();
				quest.OldIcon = this.readCell().unwrap();
				quest.OldImage = this.readCell().unwrap();
				quest.Description = this.readMultilineCell().unwrap();
				quest.Summary = this.readMultilineCell().unwrap();

				quests.push(quest);
			} catch (error) {
				Logger.error('Failed to read entry, skipping...', error);
			}
		}

		return Promise.resolve({
			result: ParsingResult.Ok,
			data: quests,
		});
	}
}
