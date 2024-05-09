import { config } from "dotenv";
import { MetadataType } from "./Metadata/MetadataType.js";
import { IDataLoader } from "./IDataLoader.js";
import { LoadQuests } from "./Quest/LoadQuests.js";
import { MetadataDb } from "./Metadata/MetadataDb.js";
import { Metadata } from "./Metadata/Metadata.js";
import { PatchDb } from "./Patches/PatchDb.js";

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

	for (let i = 0; i < patchList.length; i++) {
		const patch = patchList[i]!;
		if (patch._id.startsWith('2020-12')) {
			console.log('Reached breakpoint.');
			break;
		}

		if (meta.appliedPatches.has(patch._id)) {
			continue;
		}

		if (!loader.hasFileOfInterest(patch)) {
			continue;
		}

		console.log(`Running ${loader.name} for ${patch._id}...`);
		await loader.load(patch);

		meta.appliedPatches.add(patch._id);

		await metadataDb.updateOrCreate(meta._id, meta);
	}
}
