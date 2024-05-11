import { config } from "dotenv";

config();

export const patchesRootDir = process.env['PATCHES_DIR']!;

if (!patchesRootDir) {
	throw new Error('PATCHES_DIR is not set');
}
