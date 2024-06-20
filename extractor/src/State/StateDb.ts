import { LogRecordDao } from "../Database/LogRecordDao.js";
import { ConvertClass } from "../Utils/ConvertClass.js";
import { State } from "./DataStructures/State.js";

export class StateDb extends LogRecordDao<State> {
	constructor() {
		super('states');
	}

	protected override toInstance(data: State | null): State | null {
		if (data === null) {
			return null;
		}

		return ConvertClass.convert(data, State);
	}
}
