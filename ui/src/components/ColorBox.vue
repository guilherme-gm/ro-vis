<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	color?: number[] | null;
	defaultColor: number[];
}>();

const displayColor = computed(() => `rgb(${(props.color || props.defaultColor).join(',')})`);
</script>

<template>
	<div class="color-box" :style="{ backgroundColor: displayColor }"></div>
</template>

<style lang="scss">
	.color-box {
		display: inline-block;
		border: 2px solid black;
		width: 16px;
		height: 16px;
	}
</style>
