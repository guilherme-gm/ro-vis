import { LogRecordDao } from "../Database/LogRecordDao.js";
import { ConvertClass } from "../Utils/ConvertClass.js";
import { Quest } from "./DataStructures/Quest.js";

export class QuestDb extends LogRecordDao<Quest> {
	constructor() {
		super('quests');
	}

	protected override toInstance(data: Quest | null): Quest | null {
		if (data === null) {
			return null;
		}

		return ConvertClass.convert(data, Quest);
	}
}
