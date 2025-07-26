import { defineStore } from 'pinia';

export type Server = 'LATAM' | 'kRO';

interface ServerState {
	selectedServer: Server;
}

export const useServerStore = defineStore('server', {
	state: (): ServerState => ({
		selectedServer: 'LATAM',
	}),

	getters: {
		currentServer: (state) => state.selectedServer,
	},

	actions: {
		setServer(server: Server) {
			this.selectedServer = server;
		},
	},

	// We'll add persistence after installing the plugin
	// persist: true,
});
