import { Expose } from "class-transformer";
import { ArrayEqual } from "../../CompareUtils/ArrayEqual.js";
import { assert } from "console";

export class DescriptionLine {
	@Expose()
	public text: string = "";

	@Expose()
	public color: number[] | null = null;

	constructor(text: string, color: number[] | null) {
		this.text = text;
		this.color = color;
	}

	equals(other: DescriptionLine): boolean {
		if (this.text !== other.text) {
			return false;
		}

		if (!this.color && !other.color) {
			return true;
		}

		if ((!this.color && other.color) || (this.color && !other.color)) {
			return false;
		}

		assert(this.color && other.color);
		return ArrayEqual.isEqual(this.color!, other.color!);
	}
}
