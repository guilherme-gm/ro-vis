import { LogRecordDao } from "../Database/LogRecordDao.js";
import { Item } from "./DataStructures/Item.js";

export class ItemDb extends LogRecordDao<Item> {
	constructor() {
		super('items');
	}

	protected override toInstance(data: Item): Item {
		const i = new Item();

		i._FileVersion = data._FileVersion;
		i.Id = data.Id;
		i.IdentifiedName = data.IdentifiedName;
		i.IdentifiedDescription = data.IdentifiedDescription;
		i.IdentifiedSprite = data.IdentifiedSprite;
		i.UnidentifiedName = data.UnidentifiedName;
		i.UnidentifiedDescription = data.UnidentifiedDescription;
		i.UnidentifiedSprite = data.UnidentifiedSprite;
		i.SlotCount = data.SlotCount;
		i.IsBook = data.IsBook;
		i.CanUseBuyingStore = data.CanUseBuyingStore;
		i.CardPrefix = data.CardPrefix;
		i.CardPostfix = data.CardPostfix;
		i.CardIllustration = data.CardIllustration;

		return i;
	}
}
