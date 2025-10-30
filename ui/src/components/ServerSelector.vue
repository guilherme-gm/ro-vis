<script setup lang="ts">
import { computed } from 'vue';
import { useServerStore } from '@/stores/server';
import { ButtonVariant } from '@/components/bootstrap/Button/ButtonVariant';
import BsDropdown from '@/components/bootstrap/Dropdown/BsDropdown.vue';
import { useRoute, useRouter } from 'vue-router';

const serverStore = useServerStore();
const route = useRoute();
const router = useRouter();

const servers = [
	{ value: 'latam', label: 'LATAM' },
	{ value: 'kro', label: 'kRO Main' },
];

const selectServer = (serverId: string) => {
	serverStore.setServer(serverId as 'latam' | 'kro');
	// Navigate to the same route with updated :server param
	const name = route.name as string | undefined;
	if (name) {
		router.replace({ name, params: { ...route.params, server: serverId } });
	}
};

const currentServerLabel = computed(() => {
	const server = servers.find((s) => s.value === serverStore.currentServer);
	return server?.label ?? serverStore.currentServer;
});
</script>

<template>
	<BsDropdown :label="currentServerLabel" :variant="ButtonVariant.Secondary" :items="servers"
		@item-selected="selectServer" />
</template>
