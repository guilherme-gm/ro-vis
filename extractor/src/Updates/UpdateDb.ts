import { Filter, Sort } from "mongodb";
import { Db } from "../Database/Db.js";
import { Update } from "./Update.js";
import { plainToInstance } from "class-transformer";

export class UpdateDb extends Db<Update> {
	constructor() {
		super('updates');
	}

	public override async getAll(filter?: Filter<Update>, order?: Sort): Promise<Update[]> {
		const data = await super.getAll(filter, order);
		return data.map((d) => {
			return plainToInstance(Update, d);
		});
	}

	public override async restore(): Promise<void> {
		this.createIndex({ order: 1 });
		await super.restore();
	}
}
