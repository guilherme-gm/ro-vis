import { Expose } from "class-transformer";
import { SqlField } from "../../SqlConverter/Decorators/SqlField.js";

export class ItemMoveInfo {
	@Expose()
	@SqlField()
	public canDrop: boolean = true;

	@Expose()
	@SqlField()
	public canTrade: boolean = true;

	@Expose()
	@SqlField()
	public canMoveToStorage: boolean = true;

	@Expose()
	@SqlField()
	public canMoveToCart: boolean = true;

	@Expose()
	@SqlField()
	public canSellToNpc: boolean = true;

	@Expose()
	@SqlField()
	public canMail: boolean = true;

	@Expose()
	@SqlField()
	public canAuction: boolean = true;

	@Expose()
	@SqlField()
	public canMoveToGuildStorage: boolean = true;

	@Expose()
	@SqlField()
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
