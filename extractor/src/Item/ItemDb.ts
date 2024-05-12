import { LogRecordDao } from "../Database/LogRecordDao.js";
import { Item } from "./DataStructures/Item.js";
import { ItemMoveInfo } from "./DataStructures/ItemMoveInfo.js";

export class ItemDb extends LogRecordDao<Item> {
	constructor() {
		super('items');
	}

	protected override toInstance(data: Item | null): Item | null {
		if (data === null) {
			return null;
		}

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
		i.ClassNum = data.ClassNum;
		i.MoveInfo = new ItemMoveInfo();
		i.MoveInfo.canDrop = data.MoveInfo.canDrop;
		i.MoveInfo.canTrade = data.MoveInfo.canTrade;
		i.MoveInfo.canMoveToStorage = data.MoveInfo.canMoveToStorage;
		i.MoveInfo.canMoveToCart = data.MoveInfo.canMoveToCart;
		i.MoveInfo.canSellToNpc = data.MoveInfo.canSellToNpc;
		i.MoveInfo.canMail = data.MoveInfo.canMail;
		i.MoveInfo.canAuction = data.MoveInfo.canAuction;
		i.MoveInfo.canMoveToGuildStorage = data.MoveInfo.canMoveToGuildStorage;
		i.MoveInfo.commentName = data.MoveInfo.commentName;

		return i;
	}
}
