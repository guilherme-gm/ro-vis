import { LogRecordDao } from "../Database/LogRecordDao.js";
import { QuestV1 } from "./DataStructures/QuestV1.js";

export class QuestDb extends LogRecordDao<QuestV1> {
	constructor() {
		super('quests');
	}

	protected override toInstance(data: QuestV1): QuestV1 {
		if (data._FileVersion !== 1) {
			throw new Error('Invalid version.');
		}

		const q = new QuestV1();
		q.Description = data.Description;
		q.Id = data.Id;
		q.OldIcon = data.OldIcon;
		q.OldImage = data.OldImage;
		q.Summary = data.Summary;
		q.Title = data.Title;

		return q;
	}
}
