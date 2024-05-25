<script setup lang="ts">
import type { ItemMoveInfo } from '@/models/Item';
import ATooltip from '../ATooltip.vue';
import { useMoveInfo } from './moveInfo';

const props = defineProps<{
	from: ItemMoveInfo | null;
	to: ItemMoveInfo | null;
}>();

const { fields } = useMoveInfo();

function getClassName(key: keyof ItemMoveInfo): 'no-change' | 'added' | 'removed' {
	if (props.from?.[key] === props.to?.[key]) {
		return 'no-change';
	}

	if (props.to === null) {
		return 'removed';
	}

	return props.to[key] ? 'added' : 'removed';
}
</script>

<template>
	<div class="d-flex gap-3">
		<ATooltip v-for="(fInfo) of fields" :key="fInfo.key" :text="fInfo.label">
			<component
				:is="fInfo.icon"
				class="fs-4"
				:class="{ [getClassName(fInfo.key)]: true }"
			/>
		</ATooltip>
	</div>
</template>

<style scoped>
	.removed {
		color: rgb(236, 85, 85);
	}

	.added {
		color: rgb(89, 189, 114);
	}

	.no-change {
		color: grey;
	}
</style>
