import chalk from "chalk";
import { Logger } from "../../../Logger.js";
import { Item } from "../../DataStructures/Item.js";
import { ItemMoveInfoV5 } from "../../DataStructures/ItemMoveInfoV5.js";
import { ItemV3 } from "../../DataStructures/ItemV3.js";

export class ItemMoveInfoMerger<ItV extends ItemV3> {
	constructor(
		private itemDb: Map<number, Item>,
		private newItemMap: Map<number, ItV>,
		private moveInfoTable: ItemMoveInfoV5[] | null,
	) {}

	public mergeMoveInfo(): void {
		if (!this.moveInfoTable) {
			for (let item of this.newItemMap.values()) {
				const oldItem = this.itemDb.get(item.Id);
				if (oldItem) {
					item.MoveInfo = ItemMoveInfoV5.fromItemMoveInfo(oldItem.Id, oldItem.MoveInfo);
				} else {
					item.MoveInfo = new ItemMoveInfoV5();
				}
			}

			return;
		}

		// Resets everything, we have a new file with all entries.
		for (let item of this.newItemMap.values()) {
			item.MoveInfo = new ItemMoveInfoV5();
		}

		// Fill it
		for (let moveInfo of this.moveInfoTable) {
			const itemId = moveInfo.itemId
			const item = this.newItemMap.get(itemId);
			if (!item) {
				Logger.warn(`${chalk.whiteBright("MoveInfo")}: Item "${itemId}" is in MoveInfo but not in DB. skipping...`);
				continue;
			}

			item.MoveInfo = moveInfo;
		}
	}
}
