export class UnsupportedVersionError extends Error {
	constructor(type: string, version: number) {
		super(`Unsupported ${type} version: ${version}`);
	}
}
