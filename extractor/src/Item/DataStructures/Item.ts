import { ArrayEqual } from "../../CompareUtils/ArrayEqual.js";
import { RecordObject } from "../../Database/RecordObject.js";

/**
 * Represents a Item in the tool.
 */
export class Item implements RecordObject {
	/**
	 * The File Version that originated this object
	 */
	public _FileVersion: number = -1;

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
		);
	}
}
