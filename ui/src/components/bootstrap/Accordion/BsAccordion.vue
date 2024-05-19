<script setup lang="ts">
withDefaults(
	defineProps<{
		/**
		 * Simplifies the visualization by removing borders and rounded corners.
		 */
		flush?: boolean;
	}>(), {
		flush: false,
	},
);
</script>

<template>
	<div class="accordion" :class="{ 'accordion-flush': flush }">
		<slot />
	</div>
</template>
