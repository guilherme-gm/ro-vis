import { config } from "dotenv";
import { MetadataType } from "./Metadata/MetadataType.js";
import { IDataLoader } from "./IDataLoader.js";
import { LoadQuests } from "./Quest/LoadQuests.js";
import { MetadataDb } from "./Metadata/MetadataDb.js";
import { Metadata } from "./Metadata/Metadata.js";
import { PatchDb } from "./Patches/PatchDb.js";
import { PatchRecord } from "./Patches/PatchRecord.js";

config();

const loaders = new Map<MetadataType, IDataLoader>([
	[MetadataType.Quest, new LoadQuests()],
]);

const metadataDb = new MetadataDb();
const patchDb = new PatchDb();
const patchList = await patchDb.getAll();

for (const [metaType, loader] of loaders.entries()) {
	let meta = await metadataDb.get(metaType);
	if (meta == null) {
		meta = new Metadata(metaType);
	}

	for (let i = 0; i < 100; i++) {
		const patch = patchList[i]!;
		if (meta.appliedPatches.has(patch._id)) {
			continue;
		}

		if (!loader.hasFileOfInterest(patch)) {
			continue;
		}

		// @TODO:
		console.log(`work on ... ${patch._id}`)

		meta.appliedPatches.add(patch._id);

		await metadataDb.updateOrCreate(meta._id, meta);
	}
}
