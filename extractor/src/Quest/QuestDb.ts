import { LogRecordDao } from "../Database/LogRecordDao.js";
import { QuestV0 } from "./DataStructures/QuestV0.js";

export class QuestDb extends LogRecordDao<QuestV0> {
	constructor() {
		super('quests');
	}

	protected override toInstance(data: QuestV0): QuestV0 {
		if (data._FileVersion !== 1) {
			throw new Error('Invalid version.');
		}

		const q = new QuestV0();
		q.Description = data.Description;
		q.Id = data.Id;
		q.OldIcon = data.OldIcon;
		q.OldImage = data.OldImage;
		q.Summary = data.Summary;
		q.Title = data.Title;

		return q;
	}
}
