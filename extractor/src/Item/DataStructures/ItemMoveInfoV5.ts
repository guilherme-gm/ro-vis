import { ItemMoveInfo } from "./ItemMoveInfo.js";

export class ItemMoveInfoV5 {
	public itemId: number = 0;

	public canDrop: boolean = true;

	public canTrade: boolean = true;

	public canMoveToStorage: boolean = true;

	public canMoveToCart: boolean = true;

	public canSellToNpc: boolean = true;

	public canMail: boolean = true;

	public canAuction: boolean = true;

	public canMoveToGuildStorage: boolean = true;

	public commentName: string = "";

	public static fromItemMoveInfo(id: number, move: ItemMoveInfo): ItemMoveInfoV5 {
		const m = new ItemMoveInfoV5();
		m.itemId = id;
		m.canDrop = move.canDrop;
		m.canTrade = move.canTrade;
		m.canMoveToStorage = move.canMoveToStorage;
		m.canMoveToCart = move.canMoveToCart;
		m.canSellToNpc = move.canSellToNpc;
		m.canMail = move.canMail;
		m.canAuction = move.canAuction;
		m.canMoveToGuildStorage = move.canMoveToGuildStorage;
		m.commentName = move.commentName;

		return m;
	}

	public toEntity(): ItemMoveInfo {
		const m = new ItemMoveInfo();
		m.canDrop = this.canDrop;
		m.canTrade = this.canTrade;
		m.canMoveToStorage = this.canMoveToStorage;
		m.canMoveToCart = this.canMoveToCart;
		m.canSellToNpc = this.canSellToNpc;
		m.canMail = this.canMail;
		m.canAuction = this.canAuction;
		m.canMoveToGuildStorage = this.canMoveToGuildStorage;
		m.commentName = this.commentName;

		return m;
	}
}
