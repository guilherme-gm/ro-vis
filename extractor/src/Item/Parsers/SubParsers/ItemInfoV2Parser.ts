import { LuaTableParser } from "../../../CommonParser/LuaTableParser.js";
import { ItemV2 } from "../../DataStructures/ItemV2.js";

export type ItemInfoV2 = {
	unidentifiedDisplayName?: string;
	unidentifiedResourceName?: string;
	unidentifiedDescriptionName?: string[];
	identifiedDisplayName?: string;
	identifiedResourceName?: string;
	identifiedDescriptionName?: string[];
	slotCount?: number;
	ClassNum?: number;
};

export class ItemInfoV2Parser extends LuaTableParser<ItemInfoV2[]> {
	public static async fromFile(filePath: string): Promise<ItemInfoV2Parser> {
		return new ItemInfoV2Parser(filePath);
	}

	private readonly ExpectedKeys = new Set<string>([
		'unidentifiedDisplayName',
		'unidentifiedResourceName',
		'unidentifiedDescriptionName',
		'identifiedDisplayName',
		'identifiedResourceName',
		'identifiedDescriptionName',
		'slotCount',
		'ClassNum',
	]);

	public async parse(): Promise<ItemV2[]> {
		const items: ItemV2[] = [];

		const result = await this.extractLuaTable('tbl', true);
		Object.entries(result).forEach(([itemId, itemObj]) => {
			let additionalKeys = Object.keys(itemObj).filter((key) => !this.ExpectedKeys.has(key));
			if (additionalKeys.length > 0) {
				throw new Error(`Unexpected key found in ItemInfo V2 object. Unexpected keys: ${additionalKeys}`);
			}

			const itemV2: ItemV2 = new ItemV2();
			itemV2.Id = parseInt(itemId, 10),
				// Identified
			itemV2.IdentifiedName = itemObj.identifiedDisplayName ?? "",
			itemV2.IdentifiedDescription = itemObj.identifiedDescriptionName ?? [],
			itemV2.IdentifiedSprite = itemObj.identifiedResourceName ?? "",
				// Unidentified
			itemV2.UnidentifiedName = itemObj.unidentifiedDisplayName ?? "",
			itemV2.UnidentifiedDescription = itemObj.unidentifiedDescriptionName ?? [],
			itemV2.UnidentifiedSprite = itemObj.unidentifiedResourceName ?? "",
				// Others
			itemV2.SlotCount = itemObj.slotCount ?? 0,
			itemV2.ClassNum = itemObj.ClassNum ?? 0,

			items.push(itemV2);
		});

		return items;
	}
}
