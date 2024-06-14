import { Logger } from "../../../Logger.js";
import { KeyOfType } from "../../../Types/KeyOfType.js";
import { Item } from "../../DataStructures/Item.js";
import { ItemV } from "../../DataStructures/ItemV.js";

export class BoolIdTableMerger<ItV extends ItemV> {
	private itemDb: Map<number, Item>;

	private newItemMap: Map<number, ItV>;

	private itemClass: new () => ItV;

	constructor(
		itemDb: Map<number, Item>,
		newItemMap: Map<number, ItV>,
		itemClass: new () => ItV,
	) {
		this.itemDb = itemDb;
		this.newItemMap = newItemMap;
		this.itemClass = itemClass;
	}

	public loadBoolIdTable(
		reference: string,
		idTable: number[] | null,
		v1Key: KeyOfType<ItV, boolean>,
		itemKey: KeyOfType<Item, boolean>,
	): void {
		if (idTable) {
			for (let item of this.newItemMap.values()) {
				// @ts-ignore -- too hard to type
				item[v1Key] = false;
			}

			for (let itemId of idTable) {
				let item = this.newItemMap.get(itemId);
				if (!item) {
					item = new this.itemClass();
					item.Id = itemId;
					item.IdentifiedName = "<<Incomplete item>>";
					this.newItemMap.set(itemId, item);
					// Logger.warn(`${reference}: Item ${itemId} does not exists, could not set the flag.`);
				}

				// @ts-ignore -- too hard to type
				item[v1Key] = true;
			}
		} else {
			for (let item of this.newItemMap.values()) {
				const oldItem = this.itemDb.get(item.Id);
				if (oldItem) {
					// @ts-ignore -- too hard to type
					item[v1Key] = oldItem[itemKey];
				} else {
					// @ts-ignore -- too hard to type
					item[v1Key] = false;
				}
			}
		}
	}
}
