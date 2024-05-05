import axios from "axios";

export class Db<T> {
	private url: string;

	private apiKey: string;

	private dataSource: string;

	private database: string;

	private collection: string;

	constructor(collection: string) {
		this.url = process.env['DB_URL'] ?? '';
		this.apiKey = process.env['DB_APIKEY'] ?? '';
		this.dataSource = process.env['DB_DATASOURCE'] ?? '';
		this.database = process.env['DB_NAME'] ?? '';
		this.collection = collection;

		if (!this.url) {
			throw new Error('Missing DB_URL');
		}

		if (!this.apiKey) {
			throw new Error('Missing DB_APIKEY');
		}

		if (!this.dataSource) {
			throw new Error('Missing DB_DATASOURCE');
		}

		if (!this.database) {
			throw new Error('Missing DB_NAME');
		}
	}

	public async get(id: string): Promise<T | null> {
		const data = await axios.post(`${this.url}/action/findOne`, {
			dataSource: this.dataSource,
			database: this.database,
			collection: this.collection,
			filter: {
				_id: id,
			},
		}, {
			headers: {
				'api-key': this.apiKey,
			},
		});

		return data.data.document;
	}

	public async getAll(): Promise<T[]> {
		const response: T[] = [];

		let fetchMore = false;
		let startOffset = 0;

		do {
			const data = await axios.post(`${this.url}/action/find`, {
				dataSource: this.dataSource,
				database: this.database,
				collection: this.collection,
				skip: startOffset,
				limit: 50000,
			}, {
				headers: {
					'api-key': this.apiKey,
				},
			});

			fetchMore = data.data.documents?.length === 50000;
			startOffset += 50000;

			response.push(...data.data.documents);
		} while(fetchMore);

		return response;
	}

	public async insertMany(documents: T[]): Promise<void> {
		await axios.post(`${this.url}/action/insertMany`, {
			dataSource: this.dataSource,
			database: this.database,
			collection: this.collection,
			documents: documents,
		}, {
			headers: {
				'api-key': this.apiKey,
			},
		});
	}

	public async updateOrCreate(_id: string, document: T): Promise<void> {
		await axios.post(`${this.url}/action/updateOne`, {
			dataSource: this.dataSource,
			database: this.database,
			collection: this.collection,
			filter: {
				_id,
			},
			update: {
				$set: document,
			},
			upsert: true,
		}, {
			headers: {
				'api-key': this.apiKey,
			},
		});
	}
}
