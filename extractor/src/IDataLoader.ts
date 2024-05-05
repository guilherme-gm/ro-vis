import { PatchRecord } from "./Patches/PatchRecord.js";

export interface IDataLoader {
	hasFileOfInterest(patch: PatchRecord): boolean;

	load(patch: PatchRecord): Promise<void>;
}
