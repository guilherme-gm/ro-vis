import { CommonApi, type PatchItem } from "./CommonApi";
import type { MinQuest, Quest } from "@/models/Quest";

export type QuestPatch = PatchItem<Quest>;

export class QuestApi extends CommonApi<Quest, MinQuest> {
	public static use() {
		return new QuestApi('quests/');
	}
}
