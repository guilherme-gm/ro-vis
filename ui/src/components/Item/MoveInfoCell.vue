<script setup lang="ts">
import type { ItemMoveInfo } from '@/models/Item';
import ATooltip from '../ATooltip.vue';
import { useMoveInfo } from './moveInfo';

defineProps<{
	info: ItemMoveInfo | null;
}>();

const { fields } = useMoveInfo();
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
