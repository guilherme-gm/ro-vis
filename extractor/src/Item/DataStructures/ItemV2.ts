import { Expose } from "class-transformer";
import { IFileEntry } from "../../CommonLoader/IFileEntry.js";
import { Item } from "./Item.js";
import { ItemV } from "./ItemV.js";
import { ConvertClass } from "../../Utils/ConvertClass.js";

/**
 * Since 2012-07-11; Introduced ClassNum
 * - itemInfo.lua
 * - bookitemnametable.txt
 * - buyingstoreitemlist.txt
 * - cardpostfixnametable.txt
 * - cardprefixnametable.txt
 * - num2cardillustnametable.txt
 */
export class ItemV2 implements IFileEntry<Item> {
	public static isV2(item: ItemV): item is ItemV2 {
		return item._FileVersion === 2;
	}

	@Expose()
	public readonly _FileVersion: number = 2;

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

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public static fromItem(item: Item): ItemV2 {
		if (item._FileVersion > 2) {
			throw new Error(`Can not convert item v${item._FileVersion} to V2`);
		}

		return ConvertClass.convert(item, ItemV2);
	}

	public toEntity(): Item {
		return ConvertClass.convert(this, Item);
	}
}
