import { IFileEntry } from "../../CommonLoader/IFileEntry.js";
import { Item } from "./Item.js";
import { ItemMoveInfoV5 } from "./ItemMoveInfoV5.js";
import { ItemV } from "./ItemV.js";

/**
 * Since 2015-04-21; Introduced ItemMoveInfoV5
 * - itemInfo.lua
 * - ItemMoveInfoV5.txt
 * - bookitemnametable.txt
 * - buyingstoreitemlist.txt
 * - cardpostfixnametable.txt
 * - cardprefixnametable.txt
 * - num2cardillustnametable.txt
 */
export class ItemV3 implements IFileEntry<Item> {
	public static isV3(item: ItemV): item is ItemV3 {
		return item._FileVersion === 3;
	}

	public readonly _FileVersion: number = 3;

	/**
	 * Item ID
	 */
	public Id: number = 0;

	public IdentifiedName: string = "";

	public IdentifiedDescription: string[] = [];

	public IdentifiedSprite: string = "";

	public UnidentifiedName: string = "";

	public UnidentifiedDescription: string[] = [];

	public UnidentifiedSprite: string = "";

	public SlotCount: number = -1;

	public IsBook: boolean = false;

	public CanUseBuyingStore: boolean = false;

	public CardPrefix: string = "";

	public CardPostfix: string = "";

	public CardIllustration: string = "";

	public ClassNum: number = 0;

	public MoveInfo: ItemMoveInfoV5 = new ItemMoveInfoV5();

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public static fromItem(item: Item): ItemV3 {
		const i = new ItemV3();

		if (item._FileVersion > 3) {
			throw new Error(`Can not convert item v${item._FileVersion} to V3`);
		}

		i.Id = item.Id;
		i.IdentifiedName = item.IdentifiedName;
		i.IdentifiedDescription = item.IdentifiedDescription;
		i.IdentifiedSprite = item.IdentifiedSprite;
		i.UnidentifiedName = item.UnidentifiedName;
		i.UnidentifiedDescription = item.UnidentifiedDescription;
		i.UnidentifiedSprite = item.UnidentifiedSprite;
		i.SlotCount = item.SlotCount;
		i.IsBook = item.IsBook;
		i.CanUseBuyingStore = item.CanUseBuyingStore;
		i.CardPrefix = item.CardPrefix;
		i.CardPostfix = item.CardPostfix;
		i.CardIllustration = item.CardIllustration;
		i.ClassNum = item.ClassNum;
		i.MoveInfo = ItemMoveInfoV5.fromItemMoveInfo(item.Id, item.MoveInfo);

		return i;
	}

	public toEntity(): Item {
		const i = new Item();

		i._FileVersion = this._FileVersion;
		i.Id = this.Id;
		i.IdentifiedName = this.IdentifiedName;
		i.IdentifiedDescription = this.IdentifiedDescription;
		i.IdentifiedSprite = this.IdentifiedSprite;
		i.UnidentifiedName = this.UnidentifiedName;
		i.UnidentifiedDescription = this.UnidentifiedDescription;
		i.UnidentifiedSprite = this.UnidentifiedSprite;
		i.SlotCount = this.SlotCount;
		i.IsBook = this.IsBook;
		i.CanUseBuyingStore = this.CanUseBuyingStore;
		i.CardPrefix = this.CardPrefix;
		i.CardPostfix = this.CardPostfix;
		i.CardIllustration = this.CardIllustration;
		i.ClassNum = this.ClassNum;
		i.MoveInfo = this.MoveInfo.toEntity();

		return i;
	}
}
