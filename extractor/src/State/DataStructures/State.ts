import { Expose, Type } from "class-transformer";
import { ArrayEqual } from "../../CompareUtils/ArrayEqual.js";
import { RecordObject } from "../../Database/RecordObject.js";
import { DescriptionLine } from "./DescriptionLine.js";
import { StatePriority } from "./StatePriority.js";

/**
 * Represents a State record in the tool.
 */
export class State implements RecordObject {
	/**
	 * The File Version that originated this object
	 */
	@Expose()
	public _FileVersion: number = 1;

	/**
	 * State ID (SI / EFST_)
	 */
	@Expose()
	public Id: number = 0;

	/**
	 * EFST_ constant
	 */
	@Expose()
	public Constant: string = "";

	@Expose()
	@Type(() => DescriptionLine)
	public Description: DescriptionLine[] = [];

	@Expose()
	public HasTimeLimit: boolean = false;

	@Expose()
	public TimeLineIndex: number = -1;

	/**
	 * Whether it has an effect image (HaveEfstImgTable)
	 */
	@Expose()
	public HasEffectImage: boolean = false;

	@Expose()
	public IconImage: string = "";

	@Expose()
	public IconPriority: StatePriority = StatePriority.None;

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public equals(other: State): boolean {
		// _FileVersion is not checked, if 2 versions exists but the record is the same, we don't care.
		return (
			other.Id === this.Id
			&& other.Constant === this.Constant
			&& ArrayEqual.isEqual(other.Description, this.Description, (a, b) => a.equals(b))
			&& other.HasTimeLimit === this.HasTimeLimit
			&& other.TimeLineIndex === this.TimeLineIndex
			&& other.HasEffectImage === this.HasEffectImage
			&& other.IconImage === this.IconImage
			&& other.IconPriority === this.IconPriority
		);
	}
}
