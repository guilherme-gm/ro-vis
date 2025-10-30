<script setup lang="ts">
import { computed } from 'vue';
import { ThemeColors } from './ThemeColors';
import type { RouterLinkProps, RouteLocationNamedRaw, RouteLocationPathRaw } from 'vue-router';
import { useServerStore } from '@/stores/server';

const props = withDefaults(defineProps<{
	href?: string;
	to?: RouterLinkProps['to']
	color?: ThemeColors;
}>(), {
	color: ThemeColors.Primary,
});

const colorClass = computed(() => `link-${props.color}`);

const serverStore = useServerStore();

const resolvedTo = computed<RouterLinkProps['to'] | undefined>(() => {
	const currentServer = serverStore.currentServer;
	const t = props.to;
	if (!t) return undefined;
	if (typeof t === 'string') {
		// Prefix absolute paths
		if (t.startsWith('/')) return `/${currentServer}${t}`;
		return t;
	}
	// Object form: handle named and path variants
	const anyTo = t as any;
	if ('name' in anyTo && anyTo.name) {
		const params = { ...(anyTo.params ?? {}), server: currentServer } as Record<string, unknown>;
		return { ...anyTo, params } as RouteLocationNamedRaw;
	}
	if ('path' in anyTo && typeof anyTo.path === 'string') {
		const newPath = anyTo.path.startsWith('/') ? `/${currentServer}${anyTo.path}` : anyTo.path;
		return { ...anyTo, path: newPath } as RouteLocationPathRaw;
	}
	return t;
});
</script>

<template>
	<RouterLink v-if="resolvedTo" :to="resolvedTo" :class="{ [colorClass]: true }">
		<slot />
	</RouterLink>
	<a v-else :href="href" :class="{ [colorClass]: true }">
		<slot />
	</a>
</template>
