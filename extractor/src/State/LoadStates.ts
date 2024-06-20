import { BasicLoader } from "../CommonLoader/BasicLoader.js";
import { IDataLoader } from "../CommonLoader/IDataLoader.js";
import { UnsupportedVersionError } from "../CommonLoader/UnsupportedVersionError.js";
import { IParser } from "../CommonParser/IParser.js";
import { Logger } from "../Logger.js";
import { LogRecordSqlConverter } from "../SqlConverter/LogRecordSqlConverter.js";
import { Update } from "../Updates/Update.js";
import { State } from "./DataStructures/State.js";
import { StateV } from "./DataStructures/StateV.js";
import { StateV1Parser } from "./Parsers/StateV1Parser.js";
import { StateDb } from "./StateDb.js";

export class LoadState extends BasicLoader<State, StateV> implements IDataLoader {
	public name: string = LoadState.name;

	private readonly v1FileNames = [
		"data/luafiles514/lua files/stateicon/stateiconinfo_f.lub",
		"data/luafiles514/lua files/stateicon/efstids.lub",
		"data/luafiles514/lua files/stateicon/stateiconimginfo.lub",
		"data/luafiles514/lua files/stateicon/stateiconinfo.lub",
	];

	constructor() {
		super(new StateDb());
	}

	private getStateDataVersion(patch: Update): number {
		const date = patch._id.substring(0, 10);
		if (date.localeCompare('9999-01-01') < 0) {
			/**
			 * ??? introduces lua files for states
			 */
			return 1;
		}

		throw new Error(`State version for patch "${patch._id}" is not mapped.`);
	}

	public hasFileOfInterest(update: Update): boolean {
		const version = this.getStateDataVersion(update);
		let fileNames: string[];
		if (version === 1) {
			fileNames = this.v1FileNames;
		} else {
			throw new UnsupportedVersionError('states', version);
		}

		fileNames = fileNames.map((n) => n.toLocaleLowerCase());

		const entry = update.updates.find((f) => fileNames.includes(f.file.toLocaleLowerCase()));
		if (!entry) {
			return false;
		}

		return true;
	}

	protected async getParser(update: Update): Promise<IParser<StateV>> {
		const stateMap = new Map<number, State>();
		this.existingRecords.forEach((state) => {
			if (state.current.value !== null) {
				stateMap.set(state.current.value.Id, state.current.value);
			}
		});

		const version = this.getStateDataVersion(update);
		Logger.info(`Version: ${version}`);
		if (version === 1) {
			return new StateV1Parser(stateMap, {
				stateIconIds: this.getPathIfExists(update, 'data/luafiles514/lua files/stateicon/efstids.lub'),
				stateIconImgInfo: null, // this.getPathIfExists(update, 'data/luafiles514/lua files/stateicon/stateiconimginfo.lub'),
				stateIconInfo: null, // this.getPathIfExists(update, 'data/luafiles514/lua files/stateicon/stateiconinfo.lub'),
				stateIconFunctions: null, // this.getPathIfExists(update, 'data/luafiles514/lua files/stateicon/stateiconinfo_f.lub'),
			});
		} else {
			throw new UnsupportedVersionError('states', version);
		}
	}

	public override async dump(): Promise<void> {
		await super.dump();

		const entries = await this.entityDb.getAll();

		const sqlConverter = new LogRecordSqlConverter<State>();
		await sqlConverter.convert('states', entries);
	}
}
