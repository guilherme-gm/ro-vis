import chalk from "chalk";
import { Logger } from "../../../Logger.js";
import { KeyOfType } from "../../../Types/KeyOfType.js";
import { Item } from "../../DataStructures/Item.js";
import { ItemV } from "../../DataStructures/ItemV.js";

export class KeyValueTableMerger<ItV extends ItemV> {
	private itemDb: Map<number, Item>;

	private newItemMap: Map<number, ItV>;

	private itemClass: new () => ItV;

	constructor(
		itemDb: Map<number, Item>,
		newItemMap: Map<number, ItV>,
		itemClass: new () => ItV
	) {
		this.itemDb = itemDb;
		this.newItemMap = newItemMap;
		this.itemClass = itemClass;
	}

	public loadTable<T>(
		reference: string, table: Map<number, T> | null,
		v1Key: KeyOfType<ItV, T>,
		itemKey: KeyOfType<Item, T>,
		defaultValue?: T,
	) {
		if (v1Key === "_FileVersion") {
			throw new Error('Can not replace _FileVersion');
		}

		if (table) {
			for (let [itemId, val] of table.entries()) {
				const item = this.newItemMap.get(itemId);
				if (item) {
					// @ts-ignore -- too hard to type
					item[v1Key] = val;
				} else {
					// Logger.warn(`${chalk.whiteBright(reference)}: Item ${chalk.whiteBright(itemId)} does not exists.`);
					const newItem = new this.itemClass();
					newItem.Id = itemId;
					newItem.IdentifiedName = "<<Incomplete item>>";
					// @ts-ignore -- too hard to type
					newItem[v1Key] = val;
					this.newItemMap.set(itemId, newItem);
				}
			}
		} else {
			for (let item of this.newItemMap.values()) {
				const oldItem = this.itemDb.get(item.Id);
				if (oldItem) {
					// @ts-ignore -- too hard to type
					item[v1Key] = oldItem[itemKey];
				} else if (defaultValue !== undefined) {
					// @ts-ignore -- too hard to type
					item[v1Key] = defaultValue;
				} else {
					throw new Error(`${reference}: Item ${item.Id} is new, but not loaded.`);
				}
			}
		}
	}
}
