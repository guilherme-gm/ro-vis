<script setup lang="ts">
import { computed } from 'vue';
import { ThemeColors } from '../ThemeColors';

const props = withDefaults(
	defineProps<{
		color?: ThemeColors;
		message?: string;
	}>(),
	{
		color: ThemeColors.Primary,
		message: 'Loading...',
	},
);

const textClass = computed(() => `text-${props.color}`);
</script>

<template>
	<div class="spinner-border" role="status" :class="{ [textClass]: true }">
		<span class="visually-hidden">{{ message }}</span>
	</div>
</template>
