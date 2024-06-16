import { Exclude, Expose } from "class-transformer";
import { IFileEntry } from "../../CommonLoader/IFileEntry.js";
import { Item } from "./Item.js";
import { ItemMoveInfoV5 } from "./ItemMoveInfoV5.js";
import { ItemV } from "./ItemV.js";
import { ConvertClass } from "../../Utils/ConvertClass.js";

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

	@Expose()
	@Exclude({ toClassOnly: true })
	public readonly _FileVersion: number = 3;

	/**
	 * Item ID
	 */
	@Expose()
	public Id: number = 0;

	@Expose()
	public IdentifiedName: string = "";

	@Expose()
	public IdentifiedDescription: string[] = [];

	@Expose()
	public IdentifiedSprite: string = "";

	@Expose()
	public UnidentifiedName: string = "";

	@Expose()
	public UnidentifiedDescription: string[] = [];

	@Expose()
	public UnidentifiedSprite: string = "";

	@Expose()
	public SlotCount: number = -1;

	@Expose()
	public IsBook: boolean = false;

	@Expose()
	public CanUseBuyingStore: boolean = false;

	@Expose()
	public CardPrefix: string = "";

	@Expose()
	public CardPostfix: string = "";

	@Expose()
	public CardIllustration: string = "";

	@Expose()
	public ClassNum: number = 0;

	@Expose()
	public MoveInfo: ItemMoveInfoV5 = new ItemMoveInfoV5();

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public static fromItem(item: Item): ItemV3 {
		if (item._FileVersion > 3) {
			throw new Error(`Can not convert item v${item._FileVersion} to V3`);
		}

		const i = ConvertClass.convert(item, ItemV3);
		i.MoveInfo = ItemMoveInfoV5.fromItemMoveInfo(item.Id, item.MoveInfo);

		return i;
	}

	public toEntity(): Item {
		return ConvertClass.convert(this, Item);
	}
}
