import { Expose, Type } from "class-transformer";
import { ArrayEqual } from "../../CompareUtils/ArrayEqual.js";
import { RecordObject } from "../../Database/RecordObject.js";
import { ItemMoveInfo } from "./ItemMoveInfo.js";

/**
 * Represents a Item in the tool.
 */
export class Item implements RecordObject {
	/**
	 * The File Version that originated this object
	 */
	@Expose()
	public _FileVersion: number = -1;

	/**
	 * Item ID
	 */
	@Expose()
	public Id: number = 0;

	@Expose()
	public IdentifiedName: string = "<<Incomplete Item>>";

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
	public SlotCount: number = 0;

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
	@Type(() => ItemMoveInfo)
	public MoveInfo: ItemMoveInfo = new ItemMoveInfo();

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public equals(other: Item): boolean {
		// _FileVersion is not checked, if 2 versions exists but the record is the same, we don't care.
		return (
			this.IdentifiedName === other.IdentifiedName
			&& ArrayEqual.isEqual(this.IdentifiedDescription, other.IdentifiedDescription)
			&& this.IdentifiedSprite === other.IdentifiedSprite
			&& this.UnidentifiedName === other.UnidentifiedName
			&& ArrayEqual.isEqual(this.UnidentifiedDescription, other.UnidentifiedDescription)
			&& this.UnidentifiedSprite === other.UnidentifiedSprite
			&& this.SlotCount === other.SlotCount
			&& this.IsBook === other.IsBook
			&& this.CanUseBuyingStore === other.CanUseBuyingStore
			&& this.CardPrefix === other.CardPrefix
			&& this.CardPostfix === other.CardPostfix
			&& this.CardIllustration === other.CardIllustration
			&& this.ClassNum === other.ClassNum
			&& this.MoveInfo.equals(other.MoveInfo)
		);
	}
}
