import { Expose } from "class-transformer";

export class QuestV3RewardItem {
	@Expose()
	public ItemID: number = -1;

	@Expose()
	public ItemNum: number = -1;

	public hasChange(other: QuestV3RewardItem): boolean {
		if (!(other instanceof QuestV3RewardItem)) {
			throw new Error('Invalid type');
		}

		return (
			this.ItemID !== other.ItemID
			|| this.ItemNum !== other.ItemNum
		);
	}
}
