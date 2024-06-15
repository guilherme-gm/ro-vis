import { Expose } from "class-transformer";
import { Patch } from "./Patch.js";
import { SqlField } from "../SqlConverter/Decorators/SqlField.js";

export type UpdateItem = {
	file: string;
	patch: string;
};

export class Update {
	@Expose()
	@SqlField({
		name: 'id',
	})
	public _id: string = '';

	@Expose()
	@SqlField()
	public order: number = 0;

	@Expose()
	@SqlField({
		transform: (value: UpdateItem[]) => JSON.stringify(value),
		outType: () => String,
	})
	public updates: UpdateItem[] = [];

	@Expose()
	@SqlField({
		transform: (value: UpdateItem[]) => JSON.stringify(value),
		outType: () => String,
	})
	public patches: Patch[] = [];

	constructor(id: string) {
		this._id = id;
	}
}
