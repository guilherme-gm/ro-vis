export class QuestRewardItem {
	public ItemID: number = -1;

	public ItemNum: number = -1;

	public equals(other: QuestRewardItem): boolean {
		return (
			this.ItemID === other.ItemID
			&& this.ItemNum === other.ItemNum
		);
	}
}
