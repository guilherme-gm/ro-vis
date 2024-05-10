import {
	MongoClient,
	ServerApiVersion,
	Document,
	WithId,
	OptionalUnlessRequiredId,
	AnyBulkWriteOperation,
} from "mongodb";

export class Db<T extends Document & { _id: string; }> {
	private host: string;

	private dbUser: string;

	private dbPass: string;

	private database: string;

	private collection: string;

	private client: MongoClient;

	constructor(collection: string) {
		this.host = process.env['DB_HOST'] ?? '';
		this.dbUser = process.env['DB_USER'] ?? '';
		this.dbPass = process.env['DB_PASS'] ?? '';
		this.database = process.env['DB_NAME'] ?? '';
		this.collection = collection;

		if (!this.host) {
			throw new Error('Missing DB_URL');
		}

		if (!this.dbUser) {
			throw new Error('Missing DB_APIKEY');
		}

		if (!this.dbPass) {
			throw new Error('Missing DB_DATASOURCE');
		}

		if (!this.database) {
			throw new Error('Missing DB_NAME');
		}

		this.client = new MongoClient(`mongodb+srv://${this.dbUser}:${this.dbPass}@${this.host}/${this.database}?retryWrites=true&w=majority`, {
			serverApi: {
				version: ServerApiVersion.v1,
				strict: true,
				deprecationErrors: true,
			},
		});
	}

	public async get(id: string): Promise<T | null> {
		await this.client.connect();
		const res = await this.client.db()
			.collection<T>(this.collection)
			// @ts-expect-error -- I am not sure why, TS simply doesn't get this is string...
			.findOne({ _id: id });

		return res as T;
	}

	public async getAll(): Promise<T[]> {
		const data = await this.client.db()
			.collection<T>(this.collection)
			.find({})
			.toArray();

		return data as T[];
	}

	public async insertMany(documents: T[]): Promise<void> {
		await this.client.db()
			.collection<T>(this.collection)
			.insertMany(documents as OptionalUnlessRequiredId<T>[]);
	}

	public async updateOrCreate(_id: string, document: T): Promise<void> {
		await this.client.db()
			.collection<T>(this.collection)
			.updateOne(
				{ _id: _id } as WithId<T>,
				{ $set: document },
				{ upsert: true },
			);
	}

	public async bulkWrite(documents: T[]): Promise<void> {
		const col = this.client.db().collection<T>(this.collection);

		const changes = documents.map((doc): AnyBulkWriteOperation<T> => ({
			replaceOne: {
				// @ts-expect-error -- Weird errors...
				filter: {
					_id: doc._id,
				},
				replacement: doc,
				upsert: true,
			}
		}))

		await col.bulkWrite(changes);
	}
}
