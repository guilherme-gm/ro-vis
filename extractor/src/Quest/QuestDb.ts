import { LogRecordDao } from "../Database/LogRecordDao.js";
import { QuestV0 } from "./DataStructures/QuestV0.js";

export class QuestDb extends LogRecordDao<QuestV0> {
	constructor() {
		super('quests');
	}
}
