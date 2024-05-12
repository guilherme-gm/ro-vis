import { Expose } from "class-transformer";

export class ItemMoveInfo {
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

	public equals(other: ItemMoveInfo): boolean {
		return this.canDrop === other.canDrop
			&& this.canTrade === other.canTrade
			&& this.canMoveToStorage === other.canMoveToStorage
			&& this.canMoveToCart === other.canMoveToCart
			&& this.canSellToNpc === other.canSellToNpc
			&& this.canMail === other.canMail
			&& this.canAuction === other.canAuction
			&& this.canMoveToGuildStorage === other.canMoveToGuildStorage
			&& this.commentName === other.commentName;
	}
}
