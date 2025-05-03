/**
 * Reformats Ai4rei's patch list JSONs into a better format for extractor-go.
 *
 * The original JSON uses:
 * 1. patch name as key, but since JSON definition puts keys as unordered,
 *    it is difficult to keep the history order in extractor-go.
 *    While the input itself is not trusted to be ordered, it is the best we can have
 * 2. 1 file for ALL gpfs since 2002 and 1 file for ALL rgzs since 2002,
 *    it is hard to manually read it, and also would either make one of the types way
 *    too down the list than the other, or would require sorting in Go
 *
 * While the input order is not super trusted, since we can see cases where
 * dates are more than 1 month apart and out of order (e.g. 2011-12 being after all 2012-01-xx),
 * when both are for the same date, file order is the only info we have.
 *
 * Reformatting it with JS is MUCH easier than doing it in Go.
 *
 * The output of this script is:
 * 1. one JSON file for each year, combining RGZ and GPF patches, sorted as follows:
 *    - Extract YYYY-MM-DD date from patch name
 *    - GPF before RGZ
 *    - If 2 GPFs or 2 RGZs have the same date, input order is used
 * 2. Patch names are values of "name" key instead of being an object key
 * 3. Additional "patchDate" added to address bad naming issues,
 *    e.g. 2012-4-10 and 20120410 should be close to 2012-04-10
 * 4. Patch files that does not match the date pattern are moved to a separate "unknown.json"
 *    they are sparse and may be ignored with no issues
 *
 * This does not perform list clean up (e.g. removing sakray data)
 *
 * It is expected to have 2 JSONs in ../patches/_plist/in, named:
 * - kro-rag-rgz-hash.json
 * - kro-rag-gpf-hash.json
 *
 * Output is written to ../patches/_plist/
 */

const fs = require('fs');

function tryGetPatchDate(patchName) {
	const patterns = [
		/(\d{4})-(\d{1,2})-(\d{1,2})/, // 2024-10-30blablabla
		/(\d{4})(\d{2})(\d{2})/, // 20241030blablabla
		/(\d{4})-(\d{2})(\d{2})/, // 2024-1030blablabla
	];

	for (const pattern of patterns) {
		const dateParts = pattern.exec(patchName);
		if (!dateParts || dateParts.length !== 4) {
			continue;
		}
		const dateStr = `${dateParts[1]}-${dateParts[2].padStart(2, '0')}-${dateParts[3].padStart(2, '0')}`;
		return new Date(dateStr);
	}

	return null;
}

function getYear(patchName) {
	const year = tryGetPatchDate(patchName)?.getFullYear();
	if (!year || year < 2000 || year > 2030) {
		return 'unknown';
	}

	return year.toString();
}

function convertList(list) {
	return Object.entries(list).map(([patchName, patchData]) => {
		return {
			name: patchName,
			patchDate: tryGetPatchDate(patchName),
			files: patchData,
		};
	});
}

const gpfList = JSON.parse(fs.readFileSync('./patches/_plist/in/kro-rag-gpf-hash.json', 'utf-8'));
const rgzList = JSON.parse(fs.readFileSync('./patches/_plist/in/kro-rag-rgz-hash.json', 'utf-8'));

const allPatches = [
	...convertList(gpfList),
	...convertList(rgzList),
];

const badDatePatches = allPatches.filter((p) => getYear(p.name) === 'unknown');
const sortedPatches = allPatches
	.filter((p) => getYear(p.name) !== 'unknown')
	.sort((a, b) => a.patchDate.getTime() - b.patchDate.getTime());

const patchGroups = {
	unknown: badDatePatches,
};

for (const patch of sortedPatches) {
	const year = getYear(patch.name);
	if (!patchGroups[year]) {
		patchGroups[year] = [];
	}

	patchGroups[year].push(patch);
}

for (const patch of badDatePatches) {
	console.log(`Unknown patch date: ${patch.name}`);
}

Object.entries(patchGroups).forEach(([year, patches]) => {
	fs.writeFileSync(`./patches/_plist/${year}.json`, JSON.stringify(patches, null, 4));
});
