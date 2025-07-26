<script setup lang="ts">
	import { computed } from 'vue';
	import { useServerStore } from '@/stores/server';
	import { ButtonVariant } from '@/components/bootstrap/Button/ButtonVariant';
	import BsDropdown from '@/components/bootstrap/Dropdown/BsDropdown.vue';

	const serverStore = useServerStore();

	const servers = [
		{ value: 'LATAM', label: 'LATAM' },
		{ value: 'kRO', label: 'kRO Main' },
	];

	const selectServer = (serverId: string) => {
		serverStore.setServer(serverId as 'LATAM' | 'kRO');
	};

	const currentServerLabel = computed(() => {
		const server = servers.find((s) => s.value === serverStore.currentServer);
		return server?.label ?? serverStore.currentServer;
	});
</script>

<template>
	<BsDropdown
		:label="currentServerLabel"
		:variant="ButtonVariant.Secondary"
		:items="servers"
		@item-selected="selectServer"
	/>
</template>
