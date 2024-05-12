export class ItemMoveInfo {
	public canDrop: boolean = true;

	public canTrade: boolean = true;

	public canMoveToStorage: boolean = true;

	public canMoveToCart: boolean = true;

	public canSellToNpc: boolean = true;

	public canMail: boolean = true;

	public canAuction: boolean = true;

	public canMoveToGuildStorage: boolean = true;

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
