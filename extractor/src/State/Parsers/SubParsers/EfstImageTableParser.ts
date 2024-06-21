import { LuaTableParser } from "../../../CommonParser/LuaTableParser.js";

type HaveEfstImgTable = {
	[key: string]: number;
};

export class EfstImageTableParser extends LuaTableParser<HaveEfstImgTable> {
	public static async fromFile(efstIdPath: string, filePath: string): Promise<EfstImageTableParser> {
		return new EfstImageTableParser(efstIdPath, filePath);
	}

	private efstIdPath: string;

	constructor(efstIdPath: string, filePath: string) {
		super(filePath);

		this.efstIdPath = efstIdPath;
	}

	public async parse(): Promise<Set<number>> {
		const result = await this.extractLuaTable('HaveEfstImgTable', true, [this.efstIdPath]);

		return new Set([...Object.values(result)]);
	}
}
