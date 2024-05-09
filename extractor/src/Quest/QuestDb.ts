import { LogRecordDao } from "../Database/LogRecordDao.js";
import { Quest } from "./DataStructures/Quest.js";
import { QuestV1 } from "./DataStructures/QuestV1.js";
import { QuestV3 } from "./DataStructures/QuestV3.js";

export class QuestDb extends LogRecordDao<Quest> {
	constructor() {
		super('quests');
	}

	protected override toInstance(data: Quest): Quest {
		const version = data._FileVersion;

		if (QuestV1.isV1(data)) {
			const q = new QuestV1();
			q.Id = data.Id;
			q.Title = data.Title;
			q.Summary = data.Summary;
			q.Description = data.Description;
			q.OldIcon = data.OldIcon;
			q.OldImage = data.OldImage;

			return q;
		} else if (QuestV3.isV3(data)) {
			const q = new QuestV3();
			q.Id = data.Id;
			q.Title = data.Title;
			q.Summary = data.Summary;
			q.Description = data.Description;

			return q;
		} else {
			throw new Error(`Invalid Quest version "${version}"`);
		}
	}
}
