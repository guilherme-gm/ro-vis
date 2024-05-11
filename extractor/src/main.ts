import { config } from "dotenv";
import { MetadataType } from "./Metadata/MetadataType.js";
import { IDataLoader } from "./IDataLoader.js";
import { LoadQuests } from "./Quest/LoadQuests.js";
import { MetadataDb } from "./Metadata/MetadataDb.js";
import { Metadata } from "./Metadata/Metadata.js";
import { PatchDb } from "./Patches/PatchDb.js";
import { LoadItem } from "./Item/LoadItems.js";
import { Logger } from "./Logger.js";
import chalk from "chalk";

config();

const loaders = new Map<MetadataType, IDataLoader>([
	[MetadataType.Quest, new LoadQuests()],
	[MetadataType.Item, new LoadItem()],
]);

const metadataDb = new MetadataDb();
const patchDb = new PatchDb();
const patchList = await patchDb.getAll({}, { order: 1 });

for (const [metaType, loader] of loaders.entries()) {
	let meta = await metadataDb.get(metaType);
	if (meta == null) {
		meta = new Metadata(metaType);
	}

	for (let i = 0; i < patchList.length; i++) {
		const patch = patchList[i]!;
		if (!fs.existsSync(path.join(patchesRootDir, patch._id))) {
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
