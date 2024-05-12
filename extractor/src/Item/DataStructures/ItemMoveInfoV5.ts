import { Expose } from "class-transformer";
import { ItemMoveInfo } from "./ItemMoveInfo.js";
import { ConvertClass } from "../../Utils/ConvertClass.js";

export class ItemMoveInfoV5 {
	@Expose()
	public itemId: number = 0;

	@Expose()
	public canDrop: boolean = true;

	@Expose()
	public canTrade: boolean = true;

	@Expose()
	public canMoveToStorage: boolean = true;

	@Expose()
	public canMoveToCart: boolean = true;

	@Expose()
	public canSellToNpc: boolean = true;

	@Expose()
	public canMail: boolean = true;

	@Expose()
	public canAuction: boolean = true;

	@Expose()
	public canMoveToGuildStorage: boolean = true;

	@Expose()
	public commentName: string = "";

	public static fromItemMoveInfo(id: number, move: ItemMoveInfo): ItemMoveInfoV5 {
		const m = ConvertClass.convert(move, ItemMoveInfoV5);
		m.itemId = id;

		return m;
	}

	public toEntity(): ItemMoveInfo {
		return ConvertClass.convert(this, ItemMoveInfo);
	}
}
