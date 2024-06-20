import { LuaTableParser } from "../../../CommonParser/LuaTableParser.js";

export type EfstIdsV1 = {
	Constant: string;
	Id: number;
};

export class EfstIdsParser extends LuaTableParser<EfstIdsV1[]> {
	public static async fromFile(filePath: string): Promise<EfstIdsParser> {
		return new EfstIdsParser(filePath);
	}

	public async parse(): Promise<EfstIdsV1[]> {
		const states: EfstIdsV1[] = [];

		const result = await this.extractLuaTable('EFST_IDs', true);
		Object.entries(result).forEach(([stateConst, stateId]) => {
			// Usually there is a "__newindex" key which is a function and not a number
			if (!Number.isNaN(Number(stateId))) {
				states.push({ Constant: stateConst, Id: Number(stateId) });
			}
		});

		return states;
	}
}
