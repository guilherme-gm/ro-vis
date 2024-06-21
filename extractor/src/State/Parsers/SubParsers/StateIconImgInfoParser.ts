import { LuaTableParser } from "../../../CommonParser/LuaTableParser.js";
import { StatePriority } from '../../DataStructures/StatePriority.js';

export type StateIconImgInfoV1 = {
	EffectId: number;
	IconImage: string;
	IconPriority: number;
};

type StateIconImgList = {
	[key: number]: {
		[key: number]: string;
	};
}

export class StateIconImgInfoParser extends LuaTableParser<StateIconImgList> {
	public static async fromFile(efstIdPath: string, filePath: string): Promise<StateIconImgInfoParser> {
		return new StateIconImgInfoParser(efstIdPath, filePath);
	}

	private efstIdPath: string;

	private statePriorityMap = new Map<string, StatePriority>([
		['PRIORITY_GOLD', StatePriority.Gold],
		['PRIORITY_RED', StatePriority.Red],
		['PRIORITY_BLUE', StatePriority.Blue],
		['PRIORITY_GREEN', StatePriority.Green],
		['PRIORITY_WHITE', StatePriority.White],
	]);

	constructor(efstIdPath: string, filePath: string) {
		super(filePath);

		this.efstIdPath = efstIdPath;
	}

	public async parse(): Promise<Map<number, StateIconImgInfoV1>> {
		const states = new Map<number, StateIconImgInfoV1>();

		const result = await this.extractLuaTable('StateIconImgList', true, [this.efstIdPath]);
		Object.entries(result).forEach(([priority, stateList]) => {
			const IconPriority = this.statePriorityMap.get(priority) ?? StatePriority.None;
			if (IconPriority === StatePriority.None) {
				throw new Error(`Unknown priority: ${priority}`);
			}

			Object.entries(stateList).forEach(([effectId, icon]) => {
				states.set(Number(effectId), {
					EffectId: Number(effectId),
					IconImage: icon,
					IconPriority,
				});
			});
		});

		return states;
	}
}
