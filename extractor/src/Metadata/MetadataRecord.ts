import { Metadata } from "./Metadata.js";

export type MetadataRecord = Omit<Metadata, 'appliedPatches'> & {
	appliedPatches: string[];
};
