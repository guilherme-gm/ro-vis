export type ItemMoveInfo = {
	CanDrop: boolean;

	CanTrade: boolean;

	CanMoveToStorage: boolean;

	CanMoveToCart: boolean;

	CanSellToNpc: boolean;

	CanMail: boolean;

	CanAuction: boolean;

	CanMoveToGuildStorage: boolean;

	CommentName: string;
};

export type Item = {
	HistoryID: number;

	FileVersion: number;

	ItemID: number;

	IdentifiedName: string;

	IdentifiedDescription: string;

	IdentifiedSprite: string;

	UnidentifiedName: string;

	UnidentifiedDescription: string;

	UnidentifiedSprite: string;

	SlotCount: number;

	IsBook: boolean;

	CanUseBuyingStore: boolean;

	CardPrefix: string;

	CardIsPostfix: boolean;

	CardIllustration: string;

	ClassNum: number;

	IsCostume: boolean;

	EffectID: number;

	PackageID: number;

	MoveInfo: ItemMoveInfo;
};

export type MinItem = {
	ItemID: number;

	LastUpdate: string;

	IdentifiedName: string | null;
};
