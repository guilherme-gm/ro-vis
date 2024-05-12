import { Expose } from "class-transformer";

export class QuestRewardItem {
	@Expose()
	public ItemID: number = -1;

	@Expose()
	public ItemNum: number = -1;

	public equals(other: QuestRewardItem): boolean {
		return (
			this.ItemID === other.ItemID
			&& this.ItemNum === other.ItemNum
		);
	}
}
