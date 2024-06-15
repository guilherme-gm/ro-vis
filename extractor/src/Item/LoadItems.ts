import chalk from "chalk";
import { BasicLoader } from "../CommonLoader/BasicLoader.js";
import { IDataLoader } from "../CommonLoader/IDataLoader.js";
import { IParser } from "../CommonParser/IParser.js";
import { Logger } from "../Logger.js";
import { Item } from "./DataStructures/Item.js";
import { ItemV } from "./DataStructures/ItemV.js";
import { ItemV1Parser } from "./Parsers/ItemV1Parser.js";
import { ItemDb } from "./ItemDb.js";
import { ItemV2Parser } from "./Parsers/ItemV2Parser.js";
import { ItemV3Parser } from "./Parsers/ItemV3Parser.js";
import { Update } from "../Updates/Update.js";
import { LogRecordSqlConverter } from "../SqlConverter/LogRecordSqlConverter.js";

export class LoadItem extends BasicLoader<Item, ItemV> implements IDataLoader {
	public name: string = LoadItem.name;

	private readonly v1FileNames = [
		"data/bookitemnametable.txt",
		"data/buyingstoreitemlist.txt",
		"data/cardpostfixnametable.txt",
		"data/cardprefixnametable.txt",
		"data/idnum2itemdesctable.txt",
		"data/idnum2itemdisplaynametable.txt",
		"data/idnum2itemresnametable.txt",
		"data/itemslotcounttable.txt",
		"data/num2cardillustnametable.txt",
		"data/num2itemdesctable.txt",
		"data/num2itemdisplaynametable.txt",
		"data/num2itemresnametable.txt",
	];

	private readonly v2FileNames = [
		"System/itemInfo.lub",
		"data/bookitemnametable.txt",
		"data/buyingstoreitemlist.txt",
		"data/cardpostfixnametable.txt",
		"data/cardprefixnametable.txt",
		"data/num2cardillustnametable.txt",
	];

	private readonly v3FileNames = [
		"System/itemInfo.lub",
		"data/itemmoveinfov5.txt",
		"data/bookitemnametable.txt",
		"data/buyingstoreitemlist.txt",
		"data/cardpostfixnametable.txt",
		"data/cardprefixnametable.txt",
		"data/num2cardillustnametable.txt",
	];

	constructor() {
		super(new ItemDb());
	}

	private getItemDataVersion(update: Update): number {
		const date = update._id;
		if (date.localeCompare('2012-07-11') < 0) {
			/**
			 * Everything is txt
			 */
			return 1;
		} else if (date.localeCompare('2015-04-21') < 0) {
			// System/ItemInfo.lub
			return 2;
		} else if (date.localeCompare('2017-04-19') < 0) {
			// System/ItemInfo.lub + itemmoveinfov5
			return 3;
		} else if (date.localeCompare('9999-12-31') < 0) {
			// System/ItemInfo_true.lub
			return 4;
		}

		throw new Error(`Item version for patch "${update._id}" is not mapped.`);
	}

	public hasFileOfInterest(update: Update): boolean {
		const version = this.getItemDataVersion(update);
		let fileNames: string[];
		if (version === 1) {
			fileNames = this.v1FileNames;
		} else if (version === 2) {
			fileNames = this.v2FileNames;
		} else if (version === 3) {
			fileNames = this.v3FileNames;
		} else {
			throw new Error(`Unsupported item version "${version}"`);
		}

		fileNames = fileNames.map((n) => n.toLocaleLowerCase());

		const entry = update.updates.find((f) => fileNames.includes(f.file.toLocaleLowerCase()));
		if (!entry) {
			return false;
		}

		return true;
	}

	protected async getParser(update: Update): Promise<IParser<ItemV>> {
		const itemMap = new Map<number, Item>();
		this.existingRecords.forEach((item) => {
			if (item.current.value !== null) {
				itemMap.set(item.current.value.Id, item.current.value);
			}
		});

		const version = this.getItemDataVersion(update);
		Logger.info(`${chalk.whiteBright('Version:')} ${version}`);
		if (version === 1) {
			return new ItemV1Parser(itemMap, {
				bookItemNameTable: this.getPathIfExists(update, "data/bookitemnametable.txt"),
				buyingStoreItemList: this.getPathIfExists(update, "data/buyingstoreitemlist.txt"),
				cardPostFixNameTable: this.getPathIfExists(update, "data/cardpostfixnametable.txt"),
				cardPrefixNameTable: this.getPathIfExists(update, "data/cardprefixnametable.txt"),
				idNum2ItemDescTable: this.getPathIfExists(update, "data/idnum2itemdesctable.txt"),
				idNum2ItemDisplayNameTable: this.getPathIfExists(update, "data/idnum2itemdisplaynametable.txt"),
				idNum2ItemResNameTable: this.getPathIfExists(update, "data/idnum2itemresnametable.txt"),
				itemSlotCountTable: this.getPathIfExists(update, "data/itemslotcounttable.txt"),
				num2CardIllustNameTable: this.getPathIfExists(update, "data/num2cardillustnametable.txt"),
				num2ItemDescTable: this.getPathIfExists(update, "data/num2itemdesctable.txt"),
				num2ItemDisplayNameTable: this.getPathIfExists(update, "data/num2itemdisplaynametable.txt"),
				num2ItemResNameTable: this.getPathIfExists(update, "data/num2itemresnametable.txt"),
			});
		} else if (version === 2) {
			return new ItemV2Parser(itemMap, {
				itemInfoLua: this.getPathIfExists(update, "system/itemInfo.lub"),
				bookItemNameTable: this.getPathIfExists(update, "data/bookitemnametable.txt"),
				buyingStoreItemList: this.getPathIfExists(update, "data/buyingstoreitemlist.txt"),
				cardPostFixNameTable: this.getPathIfExists(update, "data/cardpostfixnametable.txt"),
				cardPrefixNameTable: this.getPathIfExists(update, "data/cardprefixnametable.txt"),
				num2CardIllustNameTable: this.getPathIfExists(update, "data/num2cardillustnametable.txt"),
			});
		} else if (version === 3) {
			return new ItemV3Parser(itemMap, {
				itemInfoLua: this.getPathIfExists(update, "system/itemInfo.lub"),
				moveInfoTable: this.getPathIfExists(update, "data/itemmoveinfov5.txt"),
				bookItemNameTable: this.getPathIfExists(update, "data/bookitemnametable.txt"),
				buyingStoreItemList: this.getPathIfExists(update, "data/buyingstoreitemlist.txt"),
				cardPostFixNameTable: this.getPathIfExists(update, "data/cardpostfixnametable.txt"),
				cardPrefixNameTable: this.getPathIfExists(update, "data/cardprefixnametable.txt"),
				num2CardIllustNameTable: this.getPathIfExists(update, "data/num2cardillustnametable.txt"),
			});
		} else {
			throw new Error(`Unsupported quest version "${version}"`);
		}
	}

	public override async dump(): Promise<void> {
		await super.dump();

		const entries = await this.entityDb.getAll();

		const sqlConverter = new LogRecordSqlConverter<Item>();
		await sqlConverter.convert('items', entries);
	}
}
