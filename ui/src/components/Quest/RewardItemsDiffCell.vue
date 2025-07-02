<script setup lang="ts">
import type { QuestRewardItem } from '@/models/Quest';
import DiffedValue from '../DiffedValue.vue';
import { computed } from 'vue';

const props = defineProps<{
	from: QuestRewardItem[];
	to: QuestRewardItem[];
}>();

const fromStr = computed(() => props.from.map((r) => `${r.Quantity}x ${r.ItemID}`).join('\n'));
const toStr = computed(() => props.to.map((r) => `${r.Quantity}x ${r.ItemID}`).join('\n'));
</script>

<template>
	<DiffedValue :from="fromStr" :to="toStr" />
</template>
