import {
	MongoClient,
	ServerApiVersion,
	Document,
	WithId,
	OptionalUnlessRequiredId,
	AnyBulkWriteOperation,
	Filter,
	Sort,
	IndexSpecification,
} from "mongodb";
import { Config } from "../Config/config.js";
import path from "path";
import fs from "fs";
import readline from "readline";
import { Logger } from "../Logger.js";
import { Cli } from "../Cli.js";
import { MongoServer } from "./MongoServer.js";

export class Db<T extends Document & { _id: string; }> {
	private collection: string;

	private client: MongoClient;

	constructor(collection: string) {
		this.collection = collection;

		this.client = new MongoClient(MongoServer.instance.getUri(), {
			serverApi: {
				version: ServerApiVersion.v1,
				strict: true,
				deprecationErrors: true,
			},
		});
	}

	private getDumpPath(): string {
		return path.resolve(Config.dbDir, `${this.collection}.ndjson`);
	}

	public async restore(): Promise<void> {
		if (Cli.cli.cleanRun) {
			Logger.info(`Skipping restore due to Clean Run...`);
			return;
		}

		const restoreFile = this.getDumpPath();
		if (!fs.existsSync(restoreFile)) {
			Logger.info(`Restoring "${restoreFile}"... Nothing to restore`);
			return;
		}

		Logger.info(`Restoring "${restoreFile}"...`);
		const fileStream = fs.createReadStream(restoreFile, { encoding: 'utf-8' });
		const rl = readline.createInterface({
			input: fileStream,
			crlfDelay: Infinity
		});

		let itemsToInsert: T[] = [];
		for await (const line of rl) {
			const entry = JSON.parse(line);
			itemsToInsert.push(entry);

			if (itemsToInsert.length >= 1000) {
				await this.insertMany(itemsToInsert);
				itemsToInsert = [];
			}
		}

		if (itemsToInsert.length > 0) {
			await this.insertMany(itemsToInsert);
		}

		Logger.info(`Restoring "${restoreFile}"... Done`);
	}

	public async dump(): Promise<void> {
		if (Cli.cli.dryRun) {
			Logger.info(`Dumping "${this.getDumpPath()}"... Skipped due to Dry Run`);
			return Promise.resolve();
		}

		Logger.info(`Dumping "${this.getDumpPath()}"...`);

		fs.mkdirSync(Config.dbDir, { recursive: true });

		const writeStream = fs.createWriteStream(this.getDumpPath(), { encoding: 'utf-8' });
		const items = await this.getAll();

		items.forEach((item) => {
			writeStream.write(JSON.stringify(item));
			writeStream.write('\n');
		});

		writeStream.close();

		Logger.info(`Dumping "${this.getDumpPath()}"... Done`);
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

	public async createIndex(specs: IndexSpecification): Promise<void> {
		await this.client.db()
			.collection<T>(this.collection)
			.createIndex(specs);
	}
}
