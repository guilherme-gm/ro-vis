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
	Patch: string;

	_FileVersion: number;

	Id: number;

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
	Id: string;

	Patch: string;

	IdentifiedName: string | null;
};
