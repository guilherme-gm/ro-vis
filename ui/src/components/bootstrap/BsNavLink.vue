<script setup lang="ts">
import { RouterLink, type RouterLinkProps, type RouteLocationNamedRaw, type RouteLocationPathRaw } from "vue-router";
import { computed } from 'vue';
import { useServerStore } from '@/stores/server';

const props = withDefaults(defineProps<{
	label: string;
	to?: string | RouterLinkProps["to"];
	href?: string;
	currentPage?: boolean;
}>(), {
	label: '',
	currentPage: false,
});

const serverStore = useServerStore();
const resolvedTo = computed<RouterLinkProps['to'] | undefined>(() => {
	const currentServer = serverStore.currentServer;
	const t = props.to as RouterLinkProps['to'] | undefined;
	if (!t) return undefined;
	if (typeof t === 'string') {
		if (t.startsWith('/')) return `/${currentServer}${t}`;
		return t;
	}
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
	<li class="nav-item">
		<RouterLink v-if="to" :to="(resolvedTo as any)" :aria-current="currentPage ? 'page' : 'false'" class="nav-link"
			:class="{ active: currentPage }">
			{{ label }}
		</RouterLink>
		<a v-else :href="href" :aria-current="currentPage ? 'page' : 'false'" class="nav-link"
			:class="{ active: currentPage }">
			{{ label }}
		</a>
	</li>
</template>
