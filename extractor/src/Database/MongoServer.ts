import { MongoMemoryServer } from 'mongodb-memory-server-core';
import { Cli } from "../Cli.js";
import { MemoryServerInstanceOpts } from 'mongodb-memory-server-core/lib/MongoMemoryServer.js';

export class MongoServer {
	public static instance: MongoServer = new MongoServer();

	private mongod: MongoMemoryServer | null = null;

	private uri: string = '';

	public async init(): Promise<void> {
		const opts: MemoryServerInstanceOpts = {};
		if (Cli.cli.mongoPort) {
			opts.port = Cli.cli.mongoPort;
		}

		this.mongod = await MongoMemoryServer.create({ instance: opts });

		this.uri = this.mongod.getUri('ro-vis');
	}

	public getUri(): string {
		return this.uri;
	}

	public async shutdown(): Promise<void> {
		await this.mongod?.stop();
	}
}
