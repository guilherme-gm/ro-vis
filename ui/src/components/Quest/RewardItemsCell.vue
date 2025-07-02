<script setup lang="ts">
import type { QuestRewardItem } from '@/models/Quest';

defineProps<{
	rewards: QuestRewardItem[];
}>();
</script>

<template>
	<ul v-if="rewards.length > 0">
		<li v-for="reward of rewards" :key="reward.ItemID">
			<span>{{ reward.Quantity }}x </span>
			<span>{{ reward.ItemID }}</span>
		</li>
	</ul>
	<span v-else class="text-muted">No rewards</span>
</template>
