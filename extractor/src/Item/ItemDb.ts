import { LogRecordDao } from "../Database/LogRecordDao.js";
import { Item } from "./DataStructures/Item.js";
import { ConvertClass } from "../Utils/ConvertClass.js";

export class ItemDb extends LogRecordDao<Item> {
	constructor() {
		super('items');
	}

	protected override toInstance(data: Item | null): Item | null {
		if (data === null) {
			return null;
		}

		return ConvertClass.convert(data, Item);
	}
}
