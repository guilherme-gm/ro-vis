import { MetadataType } from "./Metadata/MetadataType.js";
import { IDataLoader } from "./IDataLoader.js";
import { LoadQuests } from "./Quest/LoadQuests.js";
import { MetadataDb } from "./Metadata/MetadataDb.js";
import { Metadata } from "./Metadata/Metadata.js";
import { PatchDb } from "./Patches/PatchDb.js";
import { LoadItem } from "./Item/LoadItems.js";
import { Logger } from "./Logger.js";
import chalk from "chalk";
import * as fs from "fs";
import { Config } from "./Config/config.js";
import path from "path";
import { Cli } from "./Cli.js";
import { MongoMemoryServer } from 'mongodb-memory-server-core';
import readline from 'readline';

Cli.load();

let mongod: MongoMemoryServer | null = null;
const dryRun = Cli.cli.dryRun;
const cleanRun = Cli.cli.cleanRun;
if (dryRun) {
	mongod = await MongoMemoryServer.create();

	const uri = mongod.getUri('ro-vis');
	Config.overrideMongo(uri);

	fs.rmSync('out', { recursive: true, force: true });
	fs.mkdirSync('out');
}

let exitCode = 0;

try {
	const metadataDb = new MetadataDb();
	const patchDb = new PatchDb();
	if (dryRun) {
		await patchDb.replicate();

		if (!cleanRun) {
			await metadataDb.replicate();
		}
	}

	const patchList = await patchDb.getAll({}, { order: 1 });

	const loaders = new Map<MetadataType, IDataLoader>([
		[MetadataType.Quest, new LoadQuests()],
		[MetadataType.Item, new LoadItem()],
	]);

	for (const [metaType, loader] of loaders.entries()) {
		if (Cli.cli.only !== "" && Cli.cli.only !== metaType) {
			continue;
		}

		let meta = await metadataDb.get(metaType);
		if (meta == null) {
			meta = new Metadata(metaType);
		}

		for (let i = 0; i < patchList.length; i++) {
			const patch = patchList[i]!;

			// Skip initial patches -- for testing
			// if (patch._id.localeCompare('2018-01') < 0)
			// 	continue;

			// if (patch._id.startsWith('2020-12')) {
			// 	console.log('Reached breakpoint.');
			// 	break;
			// }

			if (!fs.existsSync(path.join(Config.patchesRootDir, patch._id))) {
				Logger.warn(`Patch ${chalk.whiteBright(patch._id)} does not exists. Skipping...`);
				continue;
			}

			if (meta.appliedPatches.has(patch._id)) {
				continue;
			}

			if (!loader.hasFileOfInterest(patch)) {
				continue;
			}

			Logger.status(`Running ${chalk.whiteBright(loader.name)} for ${chalk.white(patch._id)}...`);
			await loader.load(patch);

			meta.appliedPatches.add(patch._id);

			await metadataDb.updateOrCreate(meta._id, meta);
		}
	}
} catch (error) {
	Logger.error('An unhandled error happened...', error);
	exitCode = 1;
} finally {
	if (Cli.cli.holdProcess) {
		Logger.status('Extraction ended. Press ENTER to finish.');
		if (mongod) {
			Logger.info(`Temporary DB is running at "${chalk.whiteBright(mongod.getUri())}"`);
		}

		const readlineInterface = readline.createInterface(process.stdin, process.stdout);
		await new Promise<void>((resolve) => {
			readlineInterface.question('', () => { resolve(); });
		});
		readlineInterface.close();
	}

	if (mongod) {
		Logger.status('Closing temporary database before ending...');
		await mongod?.stop();
	}
}

process.exit(exitCode);
