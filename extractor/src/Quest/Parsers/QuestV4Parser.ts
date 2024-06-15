import * as fs from "fs";
import { LuaTableParser } from '../../CommonParser/LuaTableParser.js';
import { QuestV3RewardItem } from '../DataStructures/QuestV3RewardItem.js';
import { QuestV4 } from '../DataStructures/QuestV4.js';
import { ParsingResult } from "../../CommonParser/ParsingResult.js";
import { ParserResult } from "../../CommonParser/IParser.js";

export class QuestV4Parser extends LuaTableParser<QuestV4[]> {
	public static async fromFile(filePath: string): Promise<QuestV4Parser> {
		return new QuestV4Parser(filePath);
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
		'CoolTimeQuest',
	]);

	private readonly RewardExpectedKeys = new Set<string>(['ItemID', 'ItemNum']);

	public async parse(): Promise<ParserResult<QuestV4>> {
		if (!fs.existsSync(this.filePath)) {
			return {
				result: ParsingResult.Empty,
				data: [],
			};
		}

		const quests: QuestV4[] = [];

		const result = await this.extractLuaTable('QuestInfoList', true);
		Object.entries(result).forEach(([questId, questObj]) => {
			let additionalKeys = Object.keys(questObj).filter((key) => !this.ExpectedKeys.has(key));
			if (additionalKeys.length > 0) {
				throw new Error(`Unexpected key found in Quest V4 object. Unexpected keys: ${additionalKeys}`);
			}

			if (questObj.RewardItemList?.length > 0) {
				questObj.RewardItemList.forEach((reward) => {
					additionalKeys = Object.keys(reward).filter((key) => !this.RewardExpectedKeys.has(key));

					if (additionalKeys.length > 0) {
						throw new Error(`Unexpected key found in Quest V4 Reward object. Unexpected keys: ${additionalKeys}`);
					}
				});
			}

			const questv4 = new QuestV4();
			questv4.Id = parseInt(questId, 10);
			questv4.Title = questObj.Title ?? '';
			questv4.Description = this.fixArrayObjects(questObj.Description);
			questv4.Summary = questObj.Summary ?? '';
			questv4.IconName = questObj.IconName ?? '';
			questv4.NpcSpr = questObj.NpcSpr ?? '';
			questv4.NpcNavi = questObj.NpcNavi ?? '';
			questv4.NpcPosX = questObj.NpcPosX ?? -1;
			questv4.NpcPosY = questObj.NpcPosY ?? -1;
			questv4.RewardEXP = questObj.RewardEXP?.toString() ?? '';
			questv4.RewardJEXP = questObj.RewardJEXP?.toString() ?? '';
			questv4.CoolTimeQuest = questObj.CoolTimeQuest ?? 0;
			questv4.RewardItemList = questObj.RewardItemList?.map((reward) => {
				const r = new QuestV3RewardItem();
				r.ItemID = reward.ItemID;
				r.ItemNum = reward.ItemNum;

				return r;
			}) ?? [];

			quests.push(questv4);
		});

		return {
			result: ParsingResult.Ok,
			data: quests,
		};
	}
}
