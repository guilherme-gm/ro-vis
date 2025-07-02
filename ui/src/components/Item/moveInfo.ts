import type { ItemMoveInfo } from "@/models/Item";
import { BIconTrash, BIconEnvelopePlus, BIconCart, BIconArchive, BIconBoxes, BIconCurrencyDollar, BIconArrowLeftRight } from "bootstrap-icons-vue";
import type { Component } from "vue";
import AuctionIcon from "../icons/AuctionIcon.vue";

const fields: { key: keyof ItemMoveInfo; label: string; icon: Component; }[] = [
	{ key: "CanDrop", label: "Drop", icon: BIconTrash },
	{ key: "CanMail", label: "Mail", icon: BIconEnvelopePlus },
	{ key: "CanMoveToCart", label: "Move to cart", icon: BIconCart },
	{ key: "CanMoveToStorage", label: "Move to storage", icon: BIconArchive },
	{ key: "CanMoveToGuildStorage", label: "Move to GUILD storage", icon: BIconBoxes },
	{ key: "CanSellToNpc", label: "Sell to NPC", icon: BIconCurrencyDollar },
	{ key: "CanTrade", label: "Trade", icon: BIconArrowLeftRight },
	{ key: "CanAuction", label: "Auction", icon: AuctionIcon },
];

export function useMoveInfo() {
	return {
		fields,
	};
}
