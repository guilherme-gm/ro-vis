import { LuaTableParser } from '../../CommonParser/LuaTableParser.js';
import { QuestV3 } from "../DataStructures/QuestV3.js";

export class QuestV3Parser extends LuaTableParser<QuestV3[]> {
	public static async fromFile(filePath: string): Promise<QuestV3Parser> {
		return new QuestV3Parser(filePath);
	}

	private readonly ExpectedKeys = new Set<string>(['Title', 'Description', 'Summary'])

	public async parse(): Promise<QuestV3[]> {
		const quests: QuestV3[] = [];

		const result = await this.extractLuaTable('QuestInfoList', true);
		Object.entries(result).forEach(([questId, questObj]) => {
			const additionalKeys = Object.keys(questObj).filter((key) => !this.ExpectedKeys.has(key));
			if (additionalKeys.length > 0) {
				throw new Error(`Unexpected key found in Quest V3 object. Unexpected keys: ${additionalKeys}`);
			}

			const questv3 = new QuestV3();
			questv3.Id = parseInt(questId, 10);
			questv3.Title = questObj.Title ?? '';
			questv3.Description = questObj.Description ?? '';
			questv3.Summary = questObj.Summary ?? '';

			quests.push(questv3);
		});

		return quests;
	}
}
