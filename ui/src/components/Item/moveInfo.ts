import type { ItemMoveInfo } from "@/models/Item";
import { BIconTrash, BIconEnvelopePlus, BIconCart, BIconArchive, BIconBoxes, BIconCurrencyDollar, BIconArrowLeftRight } from "bootstrap-icons-vue";
import type { Component } from "vue";
import AuctionIcon from "../icons/AuctionIcon.vue";

const fields: { key: keyof ItemMoveInfo; label: string; icon: Component; }[] = [
	{ key: "canDrop", label: "Drop", icon: BIconTrash },
	{ key: "canMail", label: "Mail", icon: BIconEnvelopePlus },
	{ key: "canMoveToCart", label: "Move to cart", icon: BIconCart },
	{ key: "canMoveToStorage", label: "Move to storage", icon: BIconArchive },
	{ key: "canMoveToGuildStorage", label: "Move to GUILD storage", icon: BIconBoxes },
	{ key: "canSellToNpc", label: "Sell to NPC", icon: BIconCurrencyDollar },
	{ key: "canTrade", label: "Trade", icon: BIconArrowLeftRight },
	{ key: "canAuction", label: "Auction", icon: AuctionIcon },
];

export function useMoveInfo() {
	return {
		fields,
	};
}
