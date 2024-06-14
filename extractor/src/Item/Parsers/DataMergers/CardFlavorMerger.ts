import { Logger } from "../../../Logger.js";
import { Item } from "../../DataStructures/Item.js";
import { ItemV } from "../../DataStructures/ItemV.js";

export class CardFlavorMerger<ItV extends ItemV> {
	constructor(
		private itemDb: Map<number, Item>,
		private newItemMap: Map<number, ItV>,
		private cardPrefixNameTable: Map<number, string> | null = null,
		private cardPostfixNameTable: number[] | null = null,
		private itemClass: new () => ItV,
	) {}

	public loadCardFlavor(): void {
		if (!this.cardPostfixNameTable || !this.cardPrefixNameTable) {
			for (let item of this.newItemMap.values()) {
				const oldItem = this.itemDb.get(item.Id);
				if (oldItem) {
					item.CardPrefix = oldItem.CardPrefix;
					item.CardPostfix = oldItem.CardPostfix;
				} else {
					item.CardPrefix = "";
					item.CardPostfix = "";
				}
			}
		}

		// Fix items that no longer have a PostFix. In other words:
		// An item that was in postfix table but no longer is, their PostFix should now be Prefix
		if (this.cardPostfixNameTable) {
			for (let item of this.newItemMap.values()) {
				if (item.CardPostfix !== "" && !this.cardPostfixNameTable.includes(item.Id)) {
					item.CardPrefix = item.CardPostfix;
					item.CardPostfix = "";
				}
			}
		}

		// Update prefixes
		if (this.cardPrefixNameTable) {
			for (let [itemId, prefix] of this.cardPrefixNameTable) {
				let item = this.newItemMap.get(itemId);
				if (!item) {
					// Logger.warn(`${chalk.whiteBright('Card Prefix:')} Item ${itemId} doesn't exists, so we can't load its prefix.`);
					item = new this.itemClass();
					item.Id = itemId;
					item.IdentifiedName = "<<Incomplete item>>";
					this.newItemMap.set(itemId, item);
				}

				item.CardPrefix = prefix;
			}
		}

		// Update postfixes
		if (this.cardPostfixNameTable) {
			for (let itemId of this.cardPostfixNameTable) {
				let item = this.newItemMap.get(itemId);
				if (!item) {
					// Logger.error(`Card Postfix: Item ${itemId} doesn't exists, so we can't load its postfix.`);
					item = new this.itemClass();
					item.Id = itemId;
					item.IdentifiedName = "<<Incomplete item>>";
					this.newItemMap.set(itemId, item);
				}

				if (item.CardPrefix !== "") {
					item.CardPostfix = item.CardPrefix;
					item.CardPrefix = "";
				} else if (item.CardPostfix === "") {
					Logger.error(`Card Postfix: Item ${itemId} has postfix entry but does not have a prefix/postfix value.`);
				}
			}
		}
	}
}
