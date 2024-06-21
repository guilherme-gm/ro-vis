import * as fs from "fs";
import fsPromises from "fs/promises";
import os from "os";
import path from "path";
import { ParserResult } from "../../CommonParser/IParser.js";
import { ParsingResult } from "../../CommonParser/ParsingResult.js";
import { Logger } from "../../Logger.js";
import { State } from "../DataStructures/State.js";
import { StateV1 } from "../DataStructures/StateV1.js";
import { EfstIdsParser, EfstIdsV1 } from "./SubParsers/EfstIdsParser.js";
import { StateIconImgInfoParser, StateIconImgInfoV1 } from "./SubParsers/StateIconImgInfoParser.js";
import { StatePriority } from "../DataStructures/StatePriority.js";
import { StateIconInfoParser, StateInfoV1 } from "./SubParsers/StateIconInfoParser.js";
import { EfstImageTableParser } from "./SubParsers/EfstImageTableParser.js";

export type StateV1Files = {
	stateIconIds: string | null;
	stateIconFunctions: string | null;
	stateIconImgInfo: string | null;
	stateIconInfo: string | null;
};

export class StateV1Parser {
	private stateDb: Map<number, State>;

	private newStatesMap = new Map<number, StateV1>();

	private files: StateV1Files;

	private stateIconIds: EfstIdsV1[] | null = null;

	private stateImages: Map<number, StateIconImgInfoV1> | null = null;

	private stateInfo: Map<number, StateInfoV1> | null = null;

	private stateHasImageTable: Set<number> | null = null;

	constructor(stateDb: Map<number, State>, files: StateV1Files) {
		this.stateDb = stateDb;
		this.files = files;
	}

	private fileExists(path: string | null | undefined): path is string {
		if (!path) {
			return false;
		}

		const exists = fs.existsSync(path);
		if (!exists) {
			Logger.warn(`File "${path}" doesn't exists (most likely it was the same as a previous patch)`);
		}

		return exists;
	}

	private async createEfstIdsFile(): Promise<string> {
		let idMap: EfstIdsV1[];
		if (this.stateIconIds) {
			idMap = this.stateIconIds;
		} else {
			idMap = [];
			for (let state of this.stateDb.values()) {
				idMap.push({
					Id: state.Id,
					Constant: state.Constant,
				});
			}
		}

		let efstIdFile = 'EFST_IDs = {';

		for (let id of idMap) {
			efstIdFile += `\n\t${id.Constant} = ${id.Id},`;
		}

		efstIdFile += '\n}';

		const efstIdFilePath = path.join(os.tmpdir(), 'efstids.lua');
		await fsPromises.writeFile(efstIdFilePath, efstIdFile, 'utf8');

		return efstIdFilePath;
	}

	private async parseTables(): Promise<void> {
		if (this.fileExists(this.files.stateIconIds)) {
			const parser = await EfstIdsParser.fromFile(this.files.stateIconIds);
			this.stateIconIds = await parser.parse();
		}

		const efstIdsPath = await this.createEfstIdsFile();

		if (this.fileExists(this.files.stateIconImgInfo)) {
			const parser = await StateIconImgInfoParser.fromFile(efstIdsPath, this.files.stateIconImgInfo);
			this.stateImages = await parser.parse();
		}

		if (this.fileExists(this.files.stateIconInfo)) {
			const parser = await StateIconInfoParser.fromFile(efstIdsPath, this.files.stateIconInfo);
			this.stateInfo = await parser.parse();
		}

		if (this.fileExists(this.files.stateIconFunctions)) {
			const parser = await EfstImageTableParser.fromFile(efstIdsPath, this.files.stateIconFunctions);
			this.stateHasImageTable = await parser.parse();
		}
	}

	public async parse(): Promise<ParserResult<StateV1>> {
		await this.parseTables();

		this.newStatesMap = new Map<number, StateV1>();

		if (this.stateIconIds) {
			for (let parsedIds of this.stateIconIds.values()) {
				const state = new StateV1();
				state.Id = parsedIds.Id;
				state.Constant = parsedIds.Constant;

				this.newStatesMap.set(parsedIds.Id, state);
			}
		} else {
			// Assume that initially no new states are being created and we can trust
			// our db
			for (let state of this.stateDb.values()) {
				const itemId = state.Id;
				const stateV1 = StateV1.fromState(state);

				this.newStatesMap.set(itemId, stateV1);
			}
		}

		for (let state of this.newStatesMap.values()) {
			if (this.stateImages) {
				const stateImg = this.stateImages.get(state.Id);
				if (stateImg) {
					state.IconImage = stateImg.IconImage;
					state.IconPriority = stateImg.IconPriority;
					this.stateImages.delete(state.Id);
				} else {
					state.IconImage = '';
					state.IconPriority = StatePriority.None;
				}
			}

			if (this.stateInfo) {
				const stateInfo = this.stateInfo.get(state.Id);
				if (stateInfo) {
					state.Description = stateInfo.Descript;
					state.HasTimeLimit = stateInfo.HaveTimeLimit;
					state.TimeLineIndex = stateInfo.TimeLimitStrIndex;
					this.stateInfo.delete(state.Id);
				} else {
					state.Description = [];
					state.HasTimeLimit = false;
					state.TimeLineIndex = -1;
				}
			}

			if (this.stateHasImageTable) {
				state.HasEffectImage = this.stateHasImageTable.has(state.Id);
				this.stateHasImageTable.delete(state.Id);
			}
		}

		// This should never happen, since not being in EFST_ID should automatically crash the process
		// but just in case I forgot something..
		for (let stateImage of this.stateImages?.values() ?? []) {
			Logger.warn(`State "${stateImage.IconImage}" (ID: ${stateImage.EffectId}) does not exists in EFST_IDs.lua`);
		}

		// This should never happen, since not being in EFST_ID should automatically crash the process
		// but just in case I forgot something..
		for (let stateInfo of this.stateInfo?.values() ?? []) {
			Logger.warn(`State Info for ID ${stateInfo.EffectId} does not exists in EFST_IDs.lua`);
		}

		// This should never happen, since not being in EFST_ID should automatically crash the process
		// but just in case I forgot something..
		for (let stateId of this.stateHasImageTable?.values() ?? []) {
			Logger.warn(`State ID ${stateId} has image but does not exists in EFST_IDs.lua`);
		}

		return {
			result: ParsingResult.Ok,
			data: [...this.newStatesMap.values()],
		};
	}
}
