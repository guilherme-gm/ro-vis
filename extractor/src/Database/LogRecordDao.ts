import { LogRecord } from "./LogRecord.js";
import { RecordObject } from "./RecordObject.js";
import { Db } from "./Db.js";

export class LogRecordDao<T extends RecordObject> {
	private db: Db<LogRecord<T>>;

	constructor(collection: string) {
		this.db = new Db<LogRecord<T>>(collection);
	}

	public async get(id: string): Promise<LogRecord<T> | null> {
		const data = await this.db.get(id);
		if (data == null) {
			return null;
		}

		return new LogRecord(data);
	}

	public async getAll(): Promise<LogRecord<T>[]> {
		const data = await this.db.getAll();
		return data.map(d => new LogRecord(d));
	}

	public async insertMany(documents: LogRecord<T>[]): Promise<void> {
		await this.db.insertMany(documents);
	}

	public async updateOrCreate(document: LogRecord<T>): Promise<void> {
		await this.db.updateOrCreate(document._id, document);
	}
}
