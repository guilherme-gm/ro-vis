import { Expose } from "class-transformer";

export class UpdateItem {
	@Expose()
	public file: string = '';

	@Expose()
	public patch: string = '';

	constructor(file: string, patch: string) {
		this.file = file;
		this.patch = patch;
	}
}
