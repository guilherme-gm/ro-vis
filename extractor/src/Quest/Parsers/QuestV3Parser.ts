import * as fs from "fs";
import { ParserResult } from '../../CommonParser/IParser.js';
import { LuaTableParser } from '../../CommonParser/LuaTableParser.js';
import { ParsingResult } from '../../CommonParser/ParsingResult.js';
import { QuestV3 } from "../DataStructures/QuestV3.js";
import { QuestV3RewardItem } from '../DataStructures/QuestV3RewardItem.js';

export class QuestV3Parser extends LuaTableParser<QuestV3[]> {
	public static async fromFile(filePath: string): Promise<QuestV3Parser> {
		return new QuestV3Parser(filePath);
	}

	private readonly ExpectedKeys = new Set<string>([
		'Title',
		'Description',
		'Summary',
		'IconName',
		'NpcSpr',
		'NpcNavi',
		'NpcPosX',
		'NpcPosY',
		'RewardEXP',
		'RewardJEXP',
		'RewardItemList',
	]);

	private readonly RewardExpectedKeys = new Set<string>(['ItemID', 'ItemNum']);

	public async parse(): Promise<ParserResult<QuestV3>> {
		if (!fs.existsSync(this.filePath)) {
			return {
				result: ParsingResult.Empty,
				data: [],
			};
		}

		const quests: QuestV3[] = [];

		const result = await this.extractLuaTable('QuestInfoList', true);
		Object.entries(result).forEach(([questId, questObj]) => {
			let additionalKeys = Object.keys(questObj).filter((key) => !this.ExpectedKeys.has(key));
			if (additionalKeys.length > 0) {
				throw new Error(`Unexpected key found in Quest V3 object. Unexpected keys: ${additionalKeys}`);
			}

			if (questObj.RewardItemList?.length > 0) {
				questObj.RewardItemList.forEach((reward) => {
					additionalKeys = Object.keys(reward).filter((key) => !this.RewardExpectedKeys.has(key));

					if (additionalKeys.length > 0) {
						throw new Error(`Unexpected key found in Quest V3 Reward object. Unexpected keys: ${additionalKeys}`);
					}
				});
			}

			const questv3 = new QuestV3();
			questv3.Id = parseInt(questId, 10);
			questv3.Title = questObj.Title ?? '';
			questv3.Description = questObj.Description ?? '';
			questv3.Summary = questObj.Summary ?? '';
			questv3.IconName = questObj.IconName ?? '';
			questv3.NpcSpr = questObj.NpcSpr ?? '';
			questv3.NpcNavi = questObj.NpcNavi ?? '';
			questv3.NpcPosX = questObj.NpcPosX ?? -1;
			questv3.NpcPosY = questObj.NpcPosY ?? -1;
			questv3.RewardEXP = questObj.RewardEXP ?? '';
			questv3.RewardJEXP = questObj.RewardJEXP ?? '';
			questv3.RewardItemList = questObj.RewardItemList?.map((reward) => {
				const r = new QuestV3RewardItem();
				r.ItemID = reward.ItemID;
				r.ItemNum = reward.ItemNum;

				return r;
			}) ?? [];

			quests.push(questv3);
		});

		return {
			result: ParsingResult.Ok,
			data: quests,
		};
	}
}
