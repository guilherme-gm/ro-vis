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

Cli.load();

let mongod: MongoMemoryServer | null = null;
const dryRun = Cli.cli.dryRun;
const cleanRun = Cli.cli.cleanRun;
if (dryRun) {
	mongod = await MongoMemoryServer.create();

	const uri = mongod.getUri('ro-vis');
	Config.overrideMongo(uri);
}

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
		// [MetadataType.Quest, new LoadQuests()],
		[MetadataType.Item, new LoadItem()],
	]);

	for (const [metaType, loader] of loaders.entries()) {
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

			Logger.info(`Running ${chalk.whiteBright(loader.name)} for ${chalk.white(patch._id)}...`);
			await loader.load(patch);

			meta.appliedPatches.add(patch._id);

			await metadataDb.updateOrCreate(meta._id, meta);
		}
	}
} finally {
	mongod?.stop();
}
