import 'reflect-metadata';
import { MetadataType } from "./Metadata/MetadataType.js";
import { IDataLoader } from "./CommonLoader/IDataLoader.js";
import { LoadQuests } from "./Quest/LoadQuests.js";
import { MetadataDb } from "./Metadata/MetadataDb.js";
import { Metadata } from "./Metadata/Metadata.js";
import { UpdateDb } from "./Updates/UpdateDb.js";
import { LoadItem } from "./Item/LoadItems.js";
import { Logger } from "./Logger.js";
import chalk from "chalk";
import * as fs from "fs";
import { Cli } from "./Cli.js";
import readline from 'readline';
import { LoadBulkUpdateList } from "./Updates/LoadBulkPatchList.js";
import { MongoServer } from './Database/MongoServer.js';
import { SqlConverter } from './SqlConverter/SqlConverter.js';
import { UnsupportedVersionError } from './CommonLoader/UnsupportedVersionError.js';

Cli.load();

await MongoServer.instance.init();
const dryRun = Cli.cli.dryRun;
const cleanRun = Cli.cli.cleanRun;
if (dryRun || cleanRun) {
	fs.rmSync('out', { recursive: true, force: true });
	fs.mkdirSync('out');
}

let exitCode = 0;

try {
	const metadataDb = new MetadataDb();
	const updateDb = new UpdateDb();

	Logger.info('Loading MetadataDB and UpdateDB...');
	if (cleanRun) {
		const updateListLoad = new LoadBulkUpdateList();
		await updateListLoad.do();
	} else {
		await Promise.all([
			updateDb.restore(),
			metadataDb.restore(),
		]);

		const updates = await updateDb.getAll();
		if (updates.length === 0) {
			Logger.info('Loading initial update list...');

			const updateListLoad = new LoadBulkUpdateList();
			await updateListLoad.do();

			Logger.info('Initial update list loaded.');
		}
	}

	Logger.info('MetadataDB and UpdateDB loaded...');

	let patchFilter = {};
	if (Cli.cli.onlyPatches) {
		patchFilter = {
			_id: { $in: Cli.cli.onlyPatches },
		};
	}

	const updateList = await updateDb.getAll(patchFilter, { order: 1 });

	const loaders = new Map<MetadataType, IDataLoader>([
		[MetadataType.Quest, new LoadQuests()],
		[MetadataType.Item, new LoadItem()],
	]);

	for (const [metaType, loader] of loaders.entries()) {
		if (Cli.cli.only !== "" && Cli.cli.only !== metaType) {
			continue;
		}

		try {
			await loader.restore();

			let meta = await metadataDb.get(metaType);
			if (meta == null) {
				meta = new Metadata(metaType);
			}

			for (let i = 0; i < updateList.length; i++) {
				const update = updateList[i]!;

				// Skip initial patches -- for testing
				// if (patch._id.localeCompare('2018-01') < 0)
				// 	continue;

				// if (patch._id.startsWith('2020-12')) {
				// 	console.log('Reached breakpoint.');
				// 	break;
				// }

				if (meta.appliedPatches.has(update._id)) {
					continue;
				}

				if (!loader.hasFileOfInterest(update)) {
					continue;
				}

				Logger.status(`Running ${chalk.whiteBright(loader.name)} for ${chalk.white(update._id)}...`);
				await loader.load(update);

				meta.appliedPatches.add(update._id);

				await metadataDb.updateOrCreate(meta._id, meta);
			}

			await Promise.all([
				await metadataDb.dump(),
				await loader.dump(),
			]);
		} catch (error) {
			if (error instanceof UnsupportedVersionError) {
				Logger.error(error.message);
				await Promise.all([
					await metadataDb.dump(),
					await loader.dump(),
				]);
			} else {
				throw error;
			}
		}
	}

	await updateDb.dump();

	const finalUpdate = await updateDb.getAll();
	const sqlConverter = new SqlConverter();
	await sqlConverter.convert('updates', finalUpdate);
} catch (error) {
	Logger.error('An unhandled error happened...', error);
	exitCode = 1;
} finally {
	if (Cli.cli.holdProcess) {
		Logger.status('Extraction ended. Press ENTER to finish.');
		Logger.info(`Temporary DB is running at "${chalk.whiteBright(MongoServer.instance.getUri())}"`);

		const readlineInterface = readline.createInterface(process.stdin, process.stdout);
		await new Promise<void>((resolve) => {
			readlineInterface.question('', () => { resolve(); });
		});
		readlineInterface.close();
	}

	Logger.status('Closing temporary database before ending...');
	await MongoServer.instance.shutdown();
}

process.exit(exitCode);
