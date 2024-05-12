import { Expose } from "class-transformer";
import { Patch } from "./Patch.js";

export type UpdateItem = {
	file: string;
	patch: string;
};

export class Update {
	@Expose()
	public _id: string = '';

	@Expose()
	public order: number = 0;

	@Expose()
	public updates: UpdateItem[] = [];

	@Expose()
	public patches: Patch[] = [];

	constructor(id: string) {
		this._id = id;
	}
}
