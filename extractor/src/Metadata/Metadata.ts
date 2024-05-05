import { MetadataType } from "./MetadataType.js";

export class Metadata {
	public _id: MetadataType;

	public appliedPatches: Set<string> = new Set();

	constructor(type: MetadataType) {
		this._id = type;
	}
}
