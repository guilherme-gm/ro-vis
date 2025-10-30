import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { RouteName } from './RouteName'
import { useServerStore } from '@/stores/server'

const router = createRouter({
	// We can't use webHistory or GH Pages will give us 404 whenever we go to sub pages.
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/', redirect: '/kro' },

		// Scoped routes under :server
		{
			path: '/:server(latam|kro)',
			name: RouteName.Home,
			component: HomeView
		},
		{
			path: '/:server(latam|kro)/updates',
			name: RouteName.Updates,
			component: () => import("../views/UpdatesView.vue"),
		},

		// Items
		{
			path: '/:server(latam|kro)/items',
			name: RouteName.ItemList,
			component: () => import("../views/Item/ItemListView.vue"),
		},
		{
			path: '/:server(latam|kro)/items/patch/:patch',
			name: RouteName.ItemPatch,
			component: () => import("../views/Item/ItemPatchView.vue"),
		},
		{
			path: '/:server(latam|kro)/items/:itemId',
			name: RouteName.ItemHistory,
			component: () => import("../views/Item/ItemHistoryView.vue"),
		},

		// Quests
		{
			path: '/:server(latam|kro)/quests',
			name: RouteName.QuestList,
			component: () => import("../views/Quest/QuestListView.vue"),
		},
		{
			path: '/:server(latam|kro)/quests/patch/:patch',
			name: RouteName.QuestPatch,
			component: () => import("../views/Quest/QuestPatchView.vue"),
		},
		{
			path: '/:server(latam|kro)/quests/:questId',
			name: RouteName.QuestHistory,
			component: () => import("../views/Quest/QuestHistoryView.vue"),
		},

		// States
		{
			path: '/:server(latam|kro)/states',
			name: RouteName.StateList,
			component: () => import("../views/State/StateListView.vue"),
		},
		{
			path: '/:server(latam|kro)/states/patch/:patch',
			name: RouteName.StatePatch,
			component: () => import("../views/State/StatePatchView.vue"),
		},
		{
			path: '/:server(latam|kro)/states/:stateId',
			name: RouteName.StateHistory,
			component: () => import("../views/State/StateHistoryView.vue"),
		},

		// i18n
		{
			path: '/:server(latam|kro)/i18n',
			name: RouteName.I18nList,
			component: () => import("../views/I18n/I18nListView.vue"),
		},
		{
			path: '/:server(latam|kro)/i18n/patch/:patch',
			name: RouteName.I18nPatch,
			component: () => import("../views/I18n/I18nPatchView.vue"),
		},
		{
			path: '/:server(latam|kro)/i18n/:i18nId',
			name: RouteName.I18nHistory,
			component: () => import("../views/I18n/I18nHistoryView.vue"),
		},

		// Maps
		{
			path: '/:server(latam|kro)/maps',
			name: RouteName.MapList,
			component: () => import("../views/Map/MapListView.vue"),
		},
		{
			path: '/:server(latam|kro)/maps/patch/:patch',
			name: RouteName.MapPatch,
			component: () => import("../views/Map/MapPatchView.vue"),
		},
		{
			path: '/:server(latam|kro)/maps/:mapId',
			name: RouteName.MapHistory,
			component: () => import("../views/Map/MapHistoryView.vue"),
		},
	]
})

// Keep the Pinia store in sync with the URL param
router.beforeEach((to) => {
	const store = useServerStore();
	const serverParam = to.params.server as string | undefined;
	if (serverParam) {
		// Only update if changed
		if (store.currentServer !== serverParam) {
			store.setServer(serverParam as 'latam' | 'kro');
		}
	}
});

export default router
