import {
	MongoClient,
	ServerApiVersion,
	Document,
	WithId,
	OptionalUnlessRequiredId,
	AnyBulkWriteOperation,
	Filter,
	Sort,
} from "mongodb";
import { Config } from "../Config/config.js";

export class Db<T extends Document & { _id: string; }> {
	private collection: string;

	private client: MongoClient;

	constructor(collection: string) {
		const conf = Config.mainDb;
		this.collection = collection;

		let uri;
		if (conf.uri) {
			uri = conf.uri;
		} else {
			uri = `mongodb+srv://${conf.user}:${conf.pass}@${conf.host}/${conf.name}?retryWrites=true&w=majority`;
		}

		this.client = new MongoClient(uri, {
			serverApi: {
				version: ServerApiVersion.v1,
				strict: true,
				deprecationErrors: true,
			},
		});
	}

	public async replicate(): Promise<void> {
		const conf = Config.backupDb;
		const backupClient = new MongoClient(`mongodb+srv://${conf.user}:${conf.pass}@${conf.host}/${conf.name}?retryWrites=true&w=majority`, {
			serverApi: {
				version: ServerApiVersion.v1,
				strict: true,
				deprecationErrors: true,
			},
		});

		const conn = await backupClient.connect();
		const data = await conn.db().collection(this.collection).find<T>({}).toArray();
		await this.client.connect();
		await this.client.db().collection<T>(this.collection).bulkWrite(
			// @ts-ignore
			data.map((d) => ({
				replaceOne: {
					filter: { _id: d._id },
					replacement: d,
					upsert: true,
				}
			}))
		);
	}

	public async get(id: string): Promise<T | null> {
		await this.client.connect();
		const res = await this.client.db()
			.collection<T>(this.collection)
			// @ts-expect-error -- I am not sure why, TS simply doesn't get this is string...
			.findOne({ _id: id });

		return res as T;
	}

	public async getAll(filter?: Filter<T>, order?: Sort): Promise<T[]> {
		const data = await this.client.db()
			.collection<T>(this.collection)
			.find(filter ?? {})
			.sort(order ?? { _id: 1 })
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
