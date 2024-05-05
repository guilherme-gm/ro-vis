import { PatchRecord } from "./Patches/PatchRecord.js";

export interface IDataLoader {
	name: string;

	hasFileOfInterest(patch: PatchRecord): boolean;

	load(patch: PatchRecord): Promise<void>;
}
