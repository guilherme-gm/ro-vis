import { IFileEntry } from "../../CommonLoader/IFileEntry.js";
import { Item } from "./Item.js";
import { ItemV } from "./ItemV.js";
import { ConvertClass } from "../../Utils/ConvertClass.js";

/**
 * Since the beggining and maybe some additions before 2010. All TXT
 * - idnum2item<>.txt
 * - num2item<>.txt
 * - itemslotcounttable.txt
 * - bookitemnametable.txt
 * - buyingstoreitemlist.txt
 * - cardpostfixnametable.txt
 * - cardprefixnametable.txt
 * - num2cardillustnametable.txt
 */
export class ItemV1 implements IFileEntry<Item> {
	public static isV1(quest: ItemV): quest is ItemV1 {
		return quest._FileVersion === 1;
	}

	public readonly _FileVersion: number = 1;

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

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public static fromItem(item: Item): ItemV1 {
		if (item._FileVersion !== 1) {
			throw new Error(`Can not convert item v${item._FileVersion} to V1`);
		}

		return ConvertClass.convert(item, ItemV1);
	}

	public toEntity(): Item {
		return ConvertClass.convert(this, Item);
	}
}
