import * as fs from "fs";
import { PatchRecord } from "./PatchRecord.js";
// import { PatchDb } from "./PatchDb.js";

type PatchEntry = {
	name: string;
	size: number;
	hash: string;
};

type PatchList = {
	[key: string]: PatchEntry[];
}

/**
 * Process to do the initial load of patch list based on history.
 * This was a one time run to populate the db with initial data.
 */
export class LoadBulkPatchList {
	public async do(): Promise<void> {
		const gpfPatch: PatchList = JSON.parse(fs.readFileSync('./raw/kro-rag-gpf-hash.json').toString());
		const rgzPatch: PatchList = JSON.parse(fs.readFileSync('./raw/kro-rag-rgz-hash.json').toString());

		delete rgzPatch['test.rgz'];
		delete rgzPatch['gameguard_tachyonx.rgz'];
		delete rgzPatch['RO_2020-07-09_sakaRagexeRE.rgz'];
		delete gpfPatch['ExMacroData.gpf'];
		delete gpfPatch['ExMacro.gpf'];

		const gpfList = Object.entries(gpfPatch)
			.filter((entry) => entry[0].substring(0, 10).localeCompare('2012-01-01') >= 0)
			.sort((a, b) => a[0].substring(0, 10).localeCompare(b[0].substring(0, 10)));
		const rgzList = Object.entries(rgzPatch)
			.filter((entry) => entry[0].substring(0, 10).localeCompare('2012-01-01') >= 0)
			.sort((a, b) => a[0].substring(0, 10).localeCompare(b[0].substring(0, 10)));

		const mergedList: PatchRecord[] = [];
		let i = 0;
		while (gpfList.length > 0 && rgzList.length > 0) {
			let gpfHead = gpfList[0]![0].toLocaleLowerCase().substring(0, 10);
			let rgzHead = rgzList[0]![0].toLocaleLowerCase().substring(0, 10);

			// GPF comes first
			if (gpfHead.localeCompare(rgzHead) <= 0) {
				const gpfEntry = gpfList.shift()!;
				mergedList.push({
					_id: gpfEntry[0],
					order: i,
					files: gpfEntry[1]
						.filter((patch) => !!patch.hash)
						.map((patch) => patch.name),
					tags: [],
				});
			} else {
				const rgzEntry = rgzList.shift()!;
				mergedList.push({
					_id: rgzEntry[0],
					order: i,
					files: rgzEntry[1]
						.filter((patch) => !!patch.hash)
						.map((patch) => patch.name),
					tags: [],
				});
			}
			i++;
		}

		while (gpfList.length > 0) {
			const gpfEntry = gpfList.shift()!;
			mergedList.push({
				_id: gpfEntry[0],
				order: i,
				files: gpfEntry[1]
					.filter((patch) => !!patch.hash)
					.map((patch) => patch.name),
				tags: [],
			});
			i++;
		}

		while (rgzList.length > 0) {
			const rgzEntry = rgzList.shift()!;
			mergedList.push({
				_id: rgzEntry[0],
				order: i,
				files: rgzEntry[1]
					.filter((patch) => !!patch.hash)
					.map((patch) => patch.name),
				tags: [],
			});
			i++;
		}

		// Disabled to prevent accidental execution
		// const db = new PatchDb();
		// while (mergedList.length > 0) {
		// 	await db.insertMany(mergedList.splice(0, 500));
		// }
	}
}
