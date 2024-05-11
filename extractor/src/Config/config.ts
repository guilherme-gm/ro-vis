import { config } from "dotenv";

export type DbConfig = {
	uri?: string;
	host?: string;
	user?: string;
	pass?: string;
	name?: string;
};

export class Config {
	private static config: Config = new Config();

	public static get patchesRootDir(): string {
		return Config.config.patchesRootDir;
	}

	public static get mainDb(): Readonly<DbConfig> {
		return Config.config.mainDb;
	}

	public static get backupDb(): Readonly<DbConfig> {
		return Config.config.backupDb;
	}

	public static overrideMongo(uri: string): void {
		this.config.mainDb = { uri };
	}

	public patchesRootDir: string;

	public mainDb: DbConfig;

	public backupDb: DbConfig;

	constructor() {
		config();

		this.patchesRootDir = process.env['PATCHES_DIR'] ?? '';
		this.mainDb = {
			uri: process.env['DB_URI'] ?? '',
			host: process.env['DB_HOST'] ?? '',
			name: process.env['DB_NAME'] ?? '',
			user: process.env['DB_USER'] ?? '',
			pass: process.env['DB_PASS'] ?? '',
		};
		this.validateConfigs();

		this.backupDb = { ...this.mainDb };
	}

	private validateConfigs() {
		[
			{ key: 'patchesRootDir', env: 'PATCHES_DIR' },
			...(this.mainDb.uri ? [
				{ key: 'mainDb.uri', env: 'DB_URI' },
			] : [
				{ key: 'mainDb.host', env: 'DB_HOST' },
				{ key: 'mainDb.name', env: 'DB_NAME' },
				{ key: 'mainDb.user', env: 'DB_USER' },
				{ key: 'mainDb.pass', env: 'DB_PASS' },
			])
		].forEach((info) => {
			const path = info.key.split('.');
			let val = this;
			while (path.length > 0) {
				const part = path.pop();
				// @ts-ignore
				val = this[part];
			}

			if (!val) {
				throw new Error(`${info.env} is not set`);
			}
		});
	}
}
