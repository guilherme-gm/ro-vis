import { Db } from "../Database/Db.js";
import { PatchRecord } from "./PatchRecord.js";

export class PatchDb extends Db<PatchRecord> {
	constructor() {
		super('patches');
	}
}
