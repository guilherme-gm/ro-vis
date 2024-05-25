<script setup lang="ts">
import type { ItemMoveInfo } from '@/models/Item';
import { BIconArchive, BIconArrowLeftRight, BIconBoxes, BIconCart, BIconCurrencyDollar, BIconEnvelopePlus, BIconTrash } from 'bootstrap-icons-vue';
import ATooltip from '../ATooltip.vue';
import type { Component } from 'vue';
import AuctionIcon from '../icons/AuctionIcon.vue';

defineProps<{
	info: ItemMoveInfo | null;
}>();

const fields: { key: keyof ItemMoveInfo; label: string; icon: Component; }[] = [
	{ key: "canDrop", label: "Can drop", icon: BIconTrash },
	{ key: "canMail", label: "Can mail", icon: BIconEnvelopePlus },
	{ key: "canMoveToCart", label: "Can move to cart", icon: BIconCart },
	{ key: "canMoveToStorage", label: "Can move to storage", icon: BIconArchive },
	{ key: "canMoveToGuildStorage", label: "Can move to GUILD storage", icon: BIconBoxes },
	{ key: "canSellToNpc", label: "Can sell", icon: BIconCurrencyDollar },
	{ key: "canTrade", label: "Can trade", icon: BIconArrowLeftRight },
	{ key: "canAuction", label: "Can auction", icon: AuctionIcon },
];
</script>

<template>
	<span v-if="info === null">N/A</span>
	<div v-else class="d-flex gap-3">
		<ATooltip v-for="(fInfo) of fields" :key="fInfo.key" :text="fInfo.label">
			<component
				:is="fInfo.icon"
				class="fs-4"
				:class="{ disabled: !info[fInfo.key], enabled: info[fInfo.key] }"
			/>
		</ATooltip>
	</div>
</template>

<style scoped>
	.disabled {
		color: rgb(236, 85, 85);
	}

	.enabled {
		color: rgb(89, 189, 114);
	}
</style>
