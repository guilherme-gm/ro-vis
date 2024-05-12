import { Update } from "../Updates/Update.js";

export interface IDataLoader {
	name: string;

	hasFileOfInterest(patch: Update): boolean;

	load(update: Update): Promise<void>;
}
