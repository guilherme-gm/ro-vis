import * as fs from "fs";
import { ParserResult } from "../../CommonParser/IParser.js";
import { ParsingResult } from "../../CommonParser/ParsingResult.js";
import { Logger } from "../../Logger.js";
import { State } from "../DataStructures/State.js";
import { StateV1 } from "../DataStructures/StateV1.js";
import { EfstIdsParser, EfstIdsV1 } from "./SubParsers/EfstIdsParser.js";

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

	private async parseTables(): Promise<void> {
		if (this.fileExists(this.files.stateIconIds)) {
			const parser = await EfstIdsParser.fromFile(this.files.stateIconIds);
			this.stateIconIds = await parser.parse();
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

		return {
			result: ParsingResult.Ok,
			data: [...this.newStatesMap.values()],
		};
	}
}
