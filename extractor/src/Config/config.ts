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

	public static get dbDir(): string {
		return Config.config.dbDir;
	}

	public patchesRootDir: string;

	public dbDir: string;

	constructor() {
		config();

		this.patchesRootDir = process.env['PATCHES_DIR'] ?? '';
		this.dbDir = process.env['DB_DIR'] ?? '';
		this.validateConfigs();
	}

	private validateConfigs() {
		[
			{ key: 'patchesRootDir', env: 'PATCHES_DIR' },
			{ key: 'dbDir', env: 'DB_DIR' },
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
