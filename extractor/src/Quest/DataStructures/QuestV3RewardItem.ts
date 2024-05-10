export class QuestV3RewardItem {
	public ItemID: number = -1;

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
