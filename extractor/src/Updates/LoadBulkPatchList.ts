import * as fs from "fs";
import { Patch } from "./Patch.js";
import { UpdateDb } from "./UpdateDb.js";
import { Update } from "./Update.js";
import { ConvertClass } from "../Utils/ConvertClass.js";
import { UpdateItem } from "./UpdateItem.js";

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
export class LoadBulkUpdateList {
	private buildPatchList(): Patch[] {
		const gpfPatch: PatchList = JSON.parse(fs.readFileSync('./raw/kro-rag-gpf-hash.json').toString());
		const rgzPatch: PatchList = JSON.parse(fs.readFileSync('./raw/kro-rag-rgz-hash.json').toString());

		delete rgzPatch['test.rgz'];
		delete rgzPatch['gameguard_tachyonx.rgz'];
		delete rgzPatch['RO_2020-07-09_sakaRagexeRE.rgz'];
		delete gpfPatch['ExMacroData.gpf'];
		delete gpfPatch['ExMacro.gpf'];

		const gpfList = Object.entries(gpfPatch)
			.filter((entry) => entry[0].substring(0, 10).localeCompare('2012-01-01') >= 0)
			.filter((entry) => !/^\d+-\d+-\d+rData/i.test(entry[0]) && !/_sakray_/i.test(entry[0]) && !/ragexeRE/i.test(entry[0]))
			.sort((a, b) => a[0].substring(0, 10).localeCompare(b[0].substring(0, 10)));
		const rgzList = Object.entries(rgzPatch)
			.filter((entry) => entry[0].substring(0, 10).localeCompare('2012-01-01') >= 0)
			.filter((entry) => !/^\d+-\d+-\d+rData/i.test(entry[0]) && !/_sakray_/i.test(entry[0]) && !/ragexeRE/i.test(entry[0]))
			.sort((a, b) => a[0].substring(0, 10).localeCompare(b[0].substring(0, 10)));

		const patchList: Patch[] = [];

		let i = 1;
		while (gpfList.length > 0 && rgzList.length > 0) {
			let gpfHead = gpfList[0]![0].toLocaleLowerCase().substring(0, 10);
			let rgzHead = rgzList[0]![0].toLocaleLowerCase().substring(0, 10);

			// GPF comes first
			if (gpfHead.localeCompare(rgzHead) <= 0) {
				const gpfEntry = gpfList.shift()!;
				patchList.push(ConvertClass.plainToInstance(Patch, <Patch>{
					name: gpfEntry[0],
					order: i,
					files: gpfEntry[1]
						.filter((patch) => !!patch.hash)
						.map((patch) => patch.name),
				}));
			} else {
				const rgzEntry = rgzList.shift()!;
				patchList.push(ConvertClass.plainToInstance(Patch, <Patch>{
					name: rgzEntry[0],
					order: i,
					files: rgzEntry[1]
						.filter((patch) => !!patch.hash)
						.map((patch) => patch.name),
				}));
			}
			i++;
		}

		while (gpfList.length > 0) {
			const gpfEntry = gpfList.shift()!;
			patchList.push(ConvertClass.plainToInstance(Patch, <Patch>{
				name: gpfEntry[0],
				order: i,
				files: gpfEntry[1]
					.filter((patch) => !!patch.hash)
					.map((patch) => patch.name),
			}));
			i++;
		}

		while (rgzList.length > 0) {
			const rgzEntry = rgzList.shift()!;
			patchList.push(ConvertClass.plainToInstance(Patch, <Patch>{
				name: rgzEntry[0],
				order: i,
				files: rgzEntry[1]
					.filter((patch) => !!patch.hash)
					.map((patch) => patch.name),
			}));
			i++;
		}

		return patchList;
	}
	public async do(): Promise<void> {
		const patchList = this.buildPatchList();
		const updateMap = new Map<string, Update>();

		// The first patch, made up of the first presence of each of them and adapted if needed.
		// This allows us to build a full ItemDB entry that is improved as we go.
		patchList.unshift(ConvertClass.plainToInstance(Patch, <Patch>{
			"name" : "2012-01-00_dummy",
			"order" : 0,
			"files" : [
				"data\\bookitemnametable.txt",
				"data\\buyingstoreitemlist.txt",
				"data\\cardpostfixnametable.txt",
				"data\\cardprefixnametable.txt",
				"data\\idnum2itemdesctable.txt",
				"data\\idnum2itemdisplaynametable.txt",
				"data\\idnum2itemresnametable.txt",
				"data\\itemslotcounttable.txt",
				"data\\num2cardillustnametable.txt",
				"data\\num2itemdesctable.txt",
				"data\\num2itemdisplaynametable.txt",
				"data\\num2itemresnametable.txt"
			],
		}))

		patchList.sort((a, b) => a.order - b.order);

		// this is safe since patchList is sorted
		let updateOrder = 0;
		patchList.forEach((patch) => {
			const date = patch.getDate();
			let update = updateMap.get(date);
			if (!update) {
				update = new Update(date);
				update.order = updateOrder;
				updateOrder++;

				updateMap.set(date, update);
			}

			update.patches.push(patch);
		});

		for (const update of updateMap.values()) {
			const loadedFiles = new Set<string>();

			for (let i = update.patches.length - 1; i >= 0; i--) {
				const patch = update.patches[i]!;

				patch.files.forEach((file) => {
					const cleanName = file.trim().replace(/\\/g, '/');
					if (!loadedFiles.has(cleanName.toLocaleLowerCase())) {
						loadedFiles.add(cleanName.toLocaleLowerCase());
						update.updates.push(new UpdateItem(cleanName, patch.name));
					}
				});
			}
		}

		const updateList = [...updateMap.values()];
		const db = new UpdateDb();
		while (updateList.length > 0) {
			await db.insertMany(updateList.splice(0, 500));
		}
	}
}
