import { Logger } from '../../Logger.js';
import chalk  from "chalk";
import * as fs from "fs";
import { Item } from '../DataStructures/Item.js';
import { ItemV1 } from '../DataStructures/ItemV1.js';
import { BookItemNameTableV1Parser } from './SubParsers/BookItemNameTableV1Parser.js';
import { BuyingStoreItemListV1Parser } from './SubParsers/BuyingStoreItemListV1Parser.js';
import { CardIllustNameTableV1Parser } from './SubParsers/CardIllustNameTableV1Parser.js';
import { CardPostfixNameTableV1Parser } from './SubParsers/CardPostfixNameTableV1Parser.js';
import { CardPrefixNameTableV1Parser } from './SubParsers/CardPrefixNameTableV1Parser.js';
import { ItemDescTableV1Parser } from './SubParsers/ItemDescTableV1Parser.js';
import { ItemDisplayNameTableV1Parser } from './SubParsers/ItemDisplayNameTableV1Parser.js';
import { ItemResNameTableV1Parser } from './SubParsers/ItemResNameTableV1Parser.js';
import { ItemSlotCountTableV1Parser } from './SubParsers/ItemSlotCountTableV1Parser.js';
import { ParserResult } from '../../CommonParser/IParser.js';
import { ParsingResult } from '../../CommonParser/ParsingResult.js';
import { BoolIdTableMerger } from './DataMergers/BoolIdTableMerger.js';
import { KeyValueTableMerger } from './DataMergers/KeyValueTableMerger.js';
import { CardFlavorMerger } from './DataMergers/CardFlavorMerger.js';

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

	private cardIllustTable: Map<number, string> | null = null;

	private cardPrefixNameTable: Map<number, string> | null = null;

	private cardPostfixNameTable: number[] | null = null;

	constructor(itemDb: Map<number, Item>, files: ItemV1Files) {
		this.itemDb = itemDb;
		this.files = files;
	}

	private fileExists(path: string | null | undefined): path is string {
		if (!path) {
			return false;
		}

		const exists = fs.existsSync(path);
		if (!exists) {
			Logger.warn(`File "${path}" doesn't exists (most likely it was the same as a previous patch)`);
		}

		return exists;
	}

	private async parseTables(): Promise<void> {
		if (this.fileExists(this.files.idNum2ItemDisplayNameTable)) {
			const parser = await ItemDisplayNameTableV1Parser.fromFile(this.files.idNum2ItemDisplayNameTable);
			this.idNum2ItemNameTable = await parser.parse();
		}

		if (this.fileExists(this.files.idNum2ItemDescTable)) {
			const parser = await ItemDescTableV1Parser.fromFile(this.files.idNum2ItemDescTable);
			const table = await parser.parse();

			this.idNum2ItemDescTable = new Map<number, string[]>();
			for (const [key, value] of table.entries()) {
				this.idNum2ItemDescTable.set(key, value.split('\n'));
			}
		}

		if (this.fileExists(this.files.idNum2ItemResNameTable)) {
			const parser = await ItemResNameTableV1Parser.fromFile(this.files.idNum2ItemResNameTable);
			this.idNum2ItemResNameTable = await parser.parse();
		}

		if (this.fileExists(this.files.num2ItemDisplayNameTable)) {
			const parser = await ItemDisplayNameTableV1Parser.fromFile(this.files.num2ItemDisplayNameTable);
			this.num2ItemNameTable = await parser.parse();
		}

		if (this.fileExists(this.files.num2ItemDescTable)) {
			const parser = await ItemDescTableV1Parser.fromFile(this.files.num2ItemDescTable);
			const table = await parser.parse();

			this.num2ItemDescTable = new Map<number, string[]>();
			for (const [key, value] of table.entries()) {
				this.num2ItemDescTable.set(key, value.split('\n'));
			}
		}

		if (this.fileExists(this.files.num2ItemResNameTable)) {
			const parser = await ItemResNameTableV1Parser.fromFile(this.files.num2ItemResNameTable);
			this.num2ItemResNameTable = await parser.parse();
		}

		if (this.fileExists(this.files.bookItemNameTable)) {
			const parser = await BookItemNameTableV1Parser.fromFile(this.files.bookItemNameTable);
			this.bookItemNameTable = await parser.parse();
		}

		if (this.fileExists(this.files.buyingStoreItemList)) {
			const parser = await BuyingStoreItemListV1Parser.fromFile(this.files.buyingStoreItemList);
			this.buyingStoreItemList = await parser.parse();
		}

		if (this.fileExists(this.files.itemSlotCountTable)) {
			const parser = await ItemSlotCountTableV1Parser.fromFile(this.files.itemSlotCountTable);
			this.slotTable = await parser.parse();
		}

		if (this.fileExists(this.files.num2CardIllustNameTable)) {
			const parser = await CardIllustNameTableV1Parser.fromFile(this.files.num2CardIllustNameTable);
			this.cardIllustTable = await parser.parse();
		}

		if (this.fileExists(this.files.cardPostFixNameTable)) {
			const parser = await CardPostfixNameTableV1Parser.fromFile(this.files.cardPostFixNameTable);
			this.cardPostfixNameTable = await parser.parse();
		}

		if (this.fileExists(this.files.cardPrefixNameTable)) {
			const parser = await CardPrefixNameTableV1Parser.fromFile(this.files.cardPrefixNameTable);
			this.cardPrefixNameTable = await parser.parse();
		}
	}

	public async parse(): Promise<ParserResult<ItemV1>> {
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

		const keyValueTableMerger = new KeyValueTableMerger(this.itemDb, this.newItemMap, ItemV1);
		keyValueTableMerger.loadTable("Identified Item Desc Table", this.idNum2ItemDescTable, "IdentifiedDescription", "IdentifiedDescription", []);
		keyValueTableMerger.loadTable("Identified Item Res Table", this.idNum2ItemResNameTable, "IdentifiedSprite", "IdentifiedSprite", "");

		keyValueTableMerger.loadTable("Unidentified Item Name Table", this.num2ItemNameTable, "UnidentifiedName", "UnidentifiedName", "");
		keyValueTableMerger.loadTable("Unidentified Item Desc Table", this.num2ItemDescTable, "UnidentifiedDescription", "UnidentifiedDescription", []);
		keyValueTableMerger.loadTable("Unidentified Item Res Table", this.num2ItemResNameTable, "UnidentifiedSprite", "UnidentifiedSprite", "");

		keyValueTableMerger.loadTable("Slot table", this.slotTable, "SlotCount", "SlotCount", 0);
		keyValueTableMerger.loadTable("Card illust", this.cardIllustTable, "CardIllustration", "CardIllustration", "");

		const boolIdTableMerger = new BoolIdTableMerger(this.itemDb, this.newItemMap, ItemV1);
		boolIdTableMerger.loadBoolIdTable("Book", this.bookItemNameTable, "IsBook", "IsBook");
		boolIdTableMerger.loadBoolIdTable("BuyingStore", this.buyingStoreItemList, "CanUseBuyingStore", "CanUseBuyingStore");

		const cardFlavorMerger = new CardFlavorMerger(this.itemDb, this.newItemMap, this.cardPrefixNameTable, this.cardPostfixNameTable, ItemV1);
		cardFlavorMerger.loadCardFlavor();

		return {
			result: ParsingResult.Ok,
			data: [...this.newItemMap.values()],
		};
	}
}
