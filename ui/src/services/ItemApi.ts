import type { Item, MinItem } from "@/models/Item";
import { CommonApi, type PatchItem } from "./CommonApi";

export type ItemPatch = PatchItem<Item>;

export class ItemApi extends CommonApi<Item, MinItem> {
	public static use() {
		return new ItemApi('items/');
	}
}
