<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(
	defineProps<{
		/**
		 * Show as an ordered list with numbers (1, 2, 3 ...)
		 */
		numbered?: boolean;
		/**
		 * Remove borders and rounded corners. A simple line divides each item
		 */
		flush?: boolean;
	}>(), {
		numbered: false,
		flush: false,
	},
);

const element = computed(() => props.numbered ? 'ol' : 'ul');
</script>

<template>
	<component :is="element"  class="list-group" :class="{ 'list-group-numbered': numbered }">
		<slot />
	</component>
</template>
