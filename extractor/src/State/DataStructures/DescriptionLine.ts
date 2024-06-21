import { Expose } from "class-transformer";
import { ArrayEqual } from "../../CompareUtils/ArrayEqual.js";
import { assert } from "console";

export class DescriptionLine {
	@Expose()
	public Text: string = "";

	@Expose()
	public Color: number[] | null = null;

	constructor(text: string, color: number[] | null) {
		this.Text = text;
		this.Color = color;
	}

	equals(other: DescriptionLine): boolean {
		if (this.Text !== other.Text) {
			return false;
		}

		if (!this.Color && !other.Color) {
			return true;
		}

		if ((!this.Color && other.Color) || (this.Color && !other.Color)) {
			return false;
		}

		assert(this.Color && other.Color);
		return ArrayEqual.isEqual(this.Color!, other.Color!);
	}
}
