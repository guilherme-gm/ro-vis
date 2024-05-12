import { Db } from "../Database/Db.js";
import { Update } from "./Update.js";

export class UpdateDb extends Db<Update> {
	constructor() {
		super('updates');
	}
}
