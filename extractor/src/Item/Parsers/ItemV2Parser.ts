import * as fs from "fs";
import { ParserResult } from '../../CommonParser/IParser.js';
import { ParsingResult } from '../../CommonParser/ParsingResult.js';
import { Logger } from '../../Logger.js';
import { Item } from '../DataStructures/Item.js';
import { ItemV1 } from '../DataStructures/ItemV1.js';
import { ItemV2 } from '../DataStructures/ItemV2.js';
import { BoolIdTableMerger } from "./DataMergers/BoolIdTableMerger.js";
import { CardFlavorMerger } from "./DataMergers/CardFlavorMerger.js";
import { KeyValueTableMerger } from "./DataMergers/KeyValueTableMerger.js";
import { BookItemNameTableV1Parser } from './SubParsers/BookItemNameTableV1Parser.js';
import { BuyingStoreItemListV1Parser } from './SubParsers/BuyingStoreItemListV1Parser.js';
import { CardIllustNameTableV1Parser } from './SubParsers/CardIllustNameTableV1Parser.js';
import { CardPostfixNameTableV1Parser } from './SubParsers/CardPostfixNameTableV1Parser.js';
import { CardPrefixNameTableV1Parser } from './SubParsers/CardPrefixNameTableV1Parser.js';
import { ItemInfoV2Parser } from './SubParsers/ItemInfoV2Parser.js';

export type ItemV2Files = {
	itemInfoLua?: string | null;
	bookItemNameTable?: string | null;
	buyingStoreItemList?: string | null;
	cardPostFixNameTable?: string | null;
	cardPrefixNameTable?: string | null;
	num2CardIllustNameTable?: string | null;
};

export class ItemV2Parser {
	private itemDb: Map<number, Item>;

	private newItemMap = new Map<number, ItemV2>();

	private files: ItemV2Files;

	private itemInfo: ItemV2[] | null = null;

	private bookItemNameTable: number[] | null = null;

	private buyingStoreItemList: number[] | null = null;

	private cardIllustTable: Map<number, string> | null = null;

	private cardPrefixNameTable: Map<number, string> | null = null;

	private cardPostfixNameTable: number[] | null = null;

	constructor(itemDb: Map<number, Item>, files: ItemV2Files) {
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
		if (this.fileExists(this.files.itemInfoLua)) {
			const parser = await ItemInfoV2Parser.fromFile(this.files.itemInfoLua);
			this.itemInfo = await parser.parse();
		}

		if (this.fileExists(this.files.bookItemNameTable)) {
			const parser = await BookItemNameTableV1Parser.fromFile(this.files.bookItemNameTable);
			this.bookItemNameTable = await parser.parse();
		}

		if (this.fileExists(this.files.buyingStoreItemList)) {
			const parser = await BuyingStoreItemListV1Parser.fromFile(this.files.buyingStoreItemList);
			this.buyingStoreItemList = await parser.parse();
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

		this.newItemMap = new Map<number, ItemV2>();

		if (this.itemInfo) {
			for (let [itemId, itemV2] of this.itemInfo.entries()) {
				this.newItemMap.set(itemId, itemV2);
			}
		} else {
			// Assume that initially no new items are being created and we can trust
			// our db
			for (let item of this.itemDb.values()) {
				const itemId = item.Id;
				const itemV2 = ItemV2.fromItem(item);

				this.newItemMap.set(itemId, itemV2);
			}
		}

		const keyValueTableMerger = new KeyValueTableMerger(this.itemDb, this.newItemMap, ItemV2);
		keyValueTableMerger.loadTable("Card illust", this.cardIllustTable, "CardIllustration", "CardIllustration", "");

		const boolTableMerger = new BoolIdTableMerger<ItemV2>(this.itemDb, this.newItemMap, ItemV2);
		boolTableMerger.loadBoolIdTable("Book", this.bookItemNameTable, "IsBook", "IsBook");
		boolTableMerger.loadBoolIdTable("BuyingStore", this.buyingStoreItemList, "CanUseBuyingStore", "CanUseBuyingStore");

		const cardFlavorMerger = new CardFlavorMerger(this.itemDb, this.newItemMap, this.cardPrefixNameTable, this.cardPostfixNameTable, ItemV2);
		cardFlavorMerger.loadCardFlavor();

		return {
			result: ParsingResult.Ok,
			data: [...this.newItemMap.values()],
		};
	}
}
