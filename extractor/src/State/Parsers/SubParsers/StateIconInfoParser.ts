import { LuaTableParser } from "../../../CommonParser/LuaTableParser.js";
import { DescriptionLine } from "../../DataStructures/DescriptionLine.js";

export type StateInfoV1 = {
	EffectId: number;
	HaveTimeLimit: boolean;
	TimeLimitStrLine: number;
	Descript: DescriptionLine[];
};

type StateIconInfo = {
	[key: string]: {
		haveTimeLimit?: number;
		posTimeLimitStr?: number;
		descript?: [string, number[]?][];
	};
}

export class StateIconInfoParser extends LuaTableParser<StateIconInfo> {
	public static async fromFile(efstIdPath: string, filePath: string): Promise<StateIconInfoParser> {
		return new StateIconInfoParser(efstIdPath, filePath);
	}

	private efstIdPath: string;

	private expectedKeys = new Set(['haveTimeLimit', 'posTimeLimitStr', 'descript']);

	constructor(efstIdPath: string, filePath: string) {
		super(filePath);

		this.efstIdPath = efstIdPath;
	}

	public async parse(): Promise<Map<number, StateInfoV1>> {
		const states = new Map<number, StateInfoV1>();

		const result = await this.extractLuaTable('StateIconList', true, [this.efstIdPath]);
		Object.entries(result).forEach(([scId, scInfo]) => {
			Object.keys(scInfo).forEach((key) => {
				if (!this.expectedKeys.has(key)) {
					throw new Error(`Unexpected key found in StateIconInfo object. Unexpected key: ${key}`);
				}
			});

			const scIdNum = Number(scId);
			if (Number.isNaN(scIdNum)) {
				throw new Error(`StateIconInfoParser: Can not convert ${scId} to number.`);
			}

			states.set(scIdNum, {
				EffectId: scIdNum,
				HaveTimeLimit: scInfo.haveTimeLimit === 1,
				TimeLimitStrLine: (scInfo.posTimeLimitStr ?? -1),
				Descript: (scInfo.descript ?? []).map((descript) => new DescriptionLine(descript[0], descript[1] ?? [])),
			});
		});

		return states;
	}
}
