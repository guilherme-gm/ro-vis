import { Item } from '../DataStructures/Item.js';
import { ItemV1 } from '../DataStructures/ItemV1.js';
import { BookItemNameTableV1Parser } from './SubParsers/BookItemNameTableV1Parser.js';
import { BuyingStoreItemListV1Parser } from './SubParsers/BuyingStoreItemListV1Parser.js';
import { ItemDescTableV1Parser } from './SubParsers/ItemDescTableV1Parser.js';
import { ItemDisplayNameTableV1Parser } from './SubParsers/ItemDisplayNameTableV1Parser.js';
import { ItemResNameTableV1Parser } from './SubParsers/ItemResNameTableV1Parser.js';
import { ItemSlotCountTableV1Parser } from './SubParsers/ItemSlotCountTableV1Parser.js';

export type ItemV1Files = {
	bookItemNameTable?: string | null;
	buyingStoreItemList?: string | null;
	cardPostFixNameTable?: string | null;
	cardPrefixNameTable?: string | null;
	idNum2ItemDescTable?: string | null;
	idNum2ItemDisplayNameTable?: string | null;
	idNum2ItemResNameTable?: string | null;
	itemSlotCountTable?: string | null;
	num2CardIllustNameTable?: string | null;
	num2ItemDescTable?: string | null;
	num2ItemDisplayNameTable?: string | null;
	num2ItemResNameTable?: string | null;
};

type KeyOfType<T, V> = keyof {
	[P in keyof T as T[P] extends V? P: never]: any;
};

export class ItemV1Parser {
	private itemDb: Map<number, Item>;

	private newItemMap = new Map<number, ItemV1>();

	private files: ItemV1Files;

	private idNum2ItemNameTable: Map<number, string> | null = null;

	private idNum2ItemDescTable: Map<number, string[]> | null = null;

	private idNum2ItemResNameTable: Map<number, string> | null = null;

	private num2ItemNameTable: Map<number, string> | null = null;

	private num2ItemDescTable: Map<number, string[]> | null = null;

	private num2ItemResNameTable: Map<number, string> | null = null;

	private bookItemNameTable: number[] | null = null;

	private buyingStoreItemList: number[] | null = null;

	private slotTable: Map<number, number> | null = null;

	constructor(itemDb: Map<number, Item>, files: ItemV1Files) {
		this.itemDb = itemDb;
		this.files = files;
	}

	private async parseTables(): Promise<void> {
		if (this.files.idNum2ItemDisplayNameTable) {
			const parser = await ItemDisplayNameTableV1Parser.fromFile(this.files.idNum2ItemDisplayNameTable);
			this.idNum2ItemNameTable = await parser.parse();
		}

		if (this.files.idNum2ItemDescTable) {
			const parser = await ItemDescTableV1Parser.fromFile(this.files.idNum2ItemDescTable);
			const table = await parser.parse();

			this.idNum2ItemDescTable = new Map<number, string[]>();
			for (const [key, value] of table.entries()) {
				this.idNum2ItemDescTable.set(key, value.split('\n'));
			}
		}

		if (this.files.idNum2ItemResNameTable) {
			const parser = await ItemResNameTableV1Parser.fromFile(this.files.idNum2ItemResNameTable);
			this.idNum2ItemResNameTable = await parser.parse();
		}

		if (this.files.num2ItemDisplayNameTable) {
			const parser = await ItemDisplayNameTableV1Parser.fromFile(this.files.num2ItemDisplayNameTable);
			this.num2ItemNameTable = await parser.parse();
		}

		if (this.files.num2ItemDescTable) {
			const parser = await ItemDescTableV1Parser.fromFile(this.files.num2ItemDescTable);
			const table = await parser.parse();

			this.num2ItemDescTable = new Map<number, string[]>();
			for (const [key, value] of table.entries()) {
				this.num2ItemDescTable.set(key, value.split('\n'));
			}
		}

		if (this.files.num2ItemResNameTable) {
			const parser = await ItemResNameTableV1Parser.fromFile(this.files.num2ItemResNameTable);
			this.num2ItemResNameTable = await parser.parse();
		}

		if (this.files.bookItemNameTable) {
			const parser = await BookItemNameTableV1Parser.fromFile(this.files.bookItemNameTable);
			this.bookItemNameTable = await parser.parse();
		}

		if (this.files.buyingStoreItemList) {
			const parser = await BuyingStoreItemListV1Parser.fromFile(this.files.buyingStoreItemList);
			this.buyingStoreItemList = await parser.parse();
		}

		if (this.files.itemSlotCountTable) {
			const parser = await ItemSlotCountTableV1Parser.fromFile(this.files.itemSlotCountTable);
			this.slotTable = await parser.parse();
		}
	}

	private loadTable<T>(
		reference: string, table: Map<number, T> | null,
		v1Key: KeyOfType<ItemV1, T>,
		itemKey: KeyOfType<Item, T>,
	) {
		if (v1Key === "_FileVersion") {
			throw new Error('Invalid v1 key "_FileVersion"');
		}

		if (table) {
			for (let [itemId, val] of table.entries()) {
				const item = this.newItemMap.get(itemId);
				if (item) {
					// @ts-ignore -- too hard to type
					item[v1Key] = val;
				} else {
					console.error(`${reference}: Item ${itemId} does not exists.`);
				}

			}
		} else {
			for (let item of this.newItemMap.values()) {
				const oldItem = this.itemDb.get(item.Id);
				if (!oldItem) {
					throw new Error(`${reference}: Item ${item.Id} is new, but not loaded.`);
				}

				// @ts-ignore -- too hard to type
				item[v1Key] = oldItem[itemKey];
			}
		}
	}

	private loadBoolIdTable(
		reference: string,
		idTable: number[] | null,
		v1Key: KeyOfType<ItemV1, boolean>,
		itemKey: KeyOfType<Item, boolean>,
	): void {
		if (idTable) {
			for (let item of this.newItemMap.values()) {
				item[v1Key] = false;
			}

			for (let itemId of idTable) {
				const item = this.newItemMap.get(itemId);
				if (!item) {
					throw new Error(`${reference}: Item ${itemId} does not exists.`);
				}

				item[v1Key] = true;
			}
		} else {
			for (let item of this.newItemMap.values()) {
				const oldItem = this.itemDb.get(item.Id);
				if (oldItem) {
					item[v1Key] = oldItem[itemKey];
				} else {
					item[v1Key] = false;
				}
			}
		}
	}

	public async parse(): Promise<ItemV1[]> {
		await this.parseTables();

		this.newItemMap = new Map<number, ItemV1>();

		if (this.idNum2ItemNameTable) {
			for (let [itemId, itemName] of this.idNum2ItemNameTable.entries()) {
				const item = new ItemV1();
				item.Id = itemId;
				item.IdentifiedName = itemName;

				this.newItemMap.set(itemId, item);
			}
		} else {
			// Assume that initially no new items are being created and we can trust
			// our db
			for (let item of this.itemDb.values()) {
				const itemId = item.Id;
				const itemV1 = ItemV1.fromItem(item);

				this.newItemMap.set(itemId, itemV1);
			}
		}

		this.loadTable("Identified Item Desc Table", this.idNum2ItemDescTable, "IdentifiedDescription", "IdentifiedDescription");
		this.loadTable("Identified Item Res Table", this.idNum2ItemResNameTable, "IdentifiedSprite", "IdentifiedSprite");

		this.loadTable("Unidentified Item Name Table", this.num2ItemNameTable, "UnidentifiedName", "UnidentifiedName");
		this.loadTable("Unidentified Item Desc Table", this.num2ItemDescTable, "UnidentifiedDescription", "UnidentifiedDescription");
		this.loadTable("Unidentified Item Res Table", this.num2ItemResNameTable, "UnidentifiedSprite", "UnidentifiedSprite");

		this.loadTable("Slot table", this.slotTable, "SlotCount", "SlotCount");

		this.loadBoolIdTable("Book", this.bookItemNameTable, "IsBook", "IsBook");
		this.loadBoolIdTable("BuyingStore", this.buyingStoreItemList, "CanUseBuyingStore", "CanUseBuyingStore");

		return [...this.newItemMap.values()];
	}
}
