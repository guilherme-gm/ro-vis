import { Expose, Type } from "class-transformer";
import { ArrayEqual } from "../../CompareUtils/ArrayEqual.js";
import { RecordObject } from "../../Database/RecordObject.js";
import { ItemMoveInfo } from "./ItemMoveInfo.js";
import { SqlField } from "../../SqlConverter/Decorators/SqlField.js";
import { SqlNestedField } from "../../SqlConverter/Decorators/SqlNestedField.js";

/**
 * Represents a Item in the tool.
 */
export class Item implements RecordObject {
	/**
	 * The File Version that originated this object
	 */
	@Expose()
	@SqlField()
	public _FileVersion: number = -1;

	/**
	 * Item ID
	 */
	@Expose()
	@SqlField()
	public Id: number = 0;

	@Expose()
	@SqlField()
	public IdentifiedName: string = "<<Incomplete Item>>";

	@Expose()
	@SqlField({
		transform: (value: string[]) => value.join('\n'),
		outType: () => String,
	})
	public IdentifiedDescription: string[] = [];

	@Expose()
	@SqlField()
	public IdentifiedSprite: string = "";

	@Expose()
	@SqlField()
	public UnidentifiedName: string = "";

	@Expose()
	@SqlField({
		transform: (value: string[]) => value.join('\n'),
		outType: () => String,
	})
	public UnidentifiedDescription: string[] = [];

	@Expose()
	@SqlField()
	public UnidentifiedSprite: string = "";

	@Expose()
	@SqlField()
	public SlotCount: number = 0;

	@Expose()
	@SqlField()
	public IsBook: boolean = false;

	@Expose()
	@SqlField()
	public CanUseBuyingStore: boolean = false;

	@Expose()
	@SqlField()
	public CardPrefix: string = "";

	@Expose()
	@SqlField()
	public CardPostfix: string = "";

	@Expose()
	@SqlField()
	public CardIllustration: string = "";

	@Expose()
	@SqlField()
	public ClassNum: number = 0;

	@Expose()
	@Type(() => ItemMoveInfo)
	@SqlNestedField()
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
