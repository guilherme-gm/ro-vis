import { LogRecord } from "./LogRecord.js";
import { RecordObject } from "./RecordObject.js";
import { Db } from "./Db.js";

export abstract class LogRecordDao<T extends RecordObject> {
	private db: Db<LogRecord<T>>;

	constructor(collection: string) {
		this.db = new Db<LogRecord<T>>(collection);
	}

	protected abstract toInstance(data: T): T;

	private createRecord(data: LogRecord<T>): LogRecord<T> {
		const record = new LogRecord(data);

		record.current.value = this.toInstance(record.current.value);
		record.previous.forEach((prev) => {
			prev.value = this.toInstance(prev.value);
		});

		return record;
	}

	public async get(id: string): Promise<LogRecord<T> | null> {
		const data = await this.db.get(id);
		if (data == null) {
			return null;
		}

		return this.createRecord(data);
	}

	public async getAll(): Promise<LogRecord<T>[]> {
		const data = await this.db.getAll();
		return data.map(d => this.createRecord(d));
	}

	public async insertMany(documents: LogRecord<T>[]): Promise<void> {
		await this.db.insertMany(documents);
	}

	public async updateOrCreate(document: LogRecord<T>): Promise<void> {
		await this.db.updateOrCreate(document._id, document);
	}

	public async bulkWrite(documents: LogRecord<T>[]): Promise<void> {
		await this.db.bulkWrite(documents);
	}

	public async replicate(): Promise<void> {
		await this.db.replicate();
	}
}
