export type ItemMoveInfo = {
	canDrop: boolean;

	canTrade: boolean;

	canMoveToStorage: boolean;

	canMoveToCart: boolean;

	canSellToNpc: boolean;

	canMail: boolean;

	canAuction: boolean;

	canMoveToGuildStorage: boolean;

	commentName: string;
};

export type Item = {
	patch: string;

	_FileVersion: number;

	ItemId: number;

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

	CardPostfix: string;

	CardIllustration: string;

	ClassNum: number;

	MoveInfo: ItemMoveInfo;
};

export type MinItem = {
	id: string;

	patch: string;

	IdentifiedName: string | null;
};
