<script setup lang="ts">
import { computed } from 'vue';
import { ThemeColors } from './ThemeColors';
import type { RouterLinkProps } from 'vue-router';

const props = withDefaults(defineProps<{
	href?: string;
	to?: RouterLinkProps['to']
	color?: ThemeColors;
}>(), {
	color: ThemeColors.Primary,
});

const colorClass = computed(() => `link-${props.color}`);
</script>

<template>
	<RouterLink v-if="to" :to="to" :class="{ [colorClass]: true }">
		<slot />
	</RouterLink>
	<a v-else :href="href" :class="{ [colorClass]: true }">
		<slot />
	</a>
</template>
