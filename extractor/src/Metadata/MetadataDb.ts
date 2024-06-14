import { Db } from "../Database/Db.js";
import { Metadata } from "./Metadata.js";
import { MetadataRecord } from "./MetadataRecord.js";
import { MetadataType } from "./MetadataType.js";

export class MetadataDb {
	private db: Db<MetadataRecord>;

	constructor() {
		this.db = new Db('metadata');
	}

	public async get(metaType: MetadataType): Promise<Metadata | null> {
		const data = await this.db.get(metaType);
		if (data === null) {
			return data;
		}

		const obj = new Metadata(data._id);
		obj.appliedPatches = new Set(data.appliedPatches);

		return obj;
	}

	public async updateOrCreate(metaType: MetadataType, metadata: Metadata): Promise<void> {
		await this.db.updateOrCreate(metaType, {
			...metadata,
			appliedPatches: [...metadata.appliedPatches],
		});
	}

	public async dump(): Promise<void> {
		await this.db.dump();
	}

	public async restore(): Promise<void> {
		await this.db.restore();
	}
}
