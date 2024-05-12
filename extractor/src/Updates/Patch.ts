import { Expose } from "class-transformer";

export class Patch {
	@Expose()
	public name: string = '';

	@Expose()
	public order: number = 0;

	@Expose()
	public files: string[] = [];

	getDate(): string {
		let dateMatch = this.name.match(/^(\d{4})-(\d{1,2})-(\d{1,2})/);
		if (dateMatch) {
			return `${dateMatch[1]}-${dateMatch[2]?.padStart(2, '0')}-${dateMatch[3]?.padStart(2, '0')}`;
		}

		dateMatch = this.name.match(/^(\d{4})-(\d{4})/);
		if (dateMatch) {
			return `${dateMatch[1]}-${dateMatch[2]?.substring(0, 2)?.padStart(2, '0')}-${dateMatch[2]?.substring(2, 4)?.padStart(2, '0')}`;
		}

		throw new Error(`Can not extract date from patch "${this.name}".`);
	}
}
