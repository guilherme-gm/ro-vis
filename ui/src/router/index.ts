import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { RouteName } from './RouteName'

const router = createRouter({
	// We can't use webHistory or GH Pages will give us 404 whenever we go to sub pages.
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			name: RouteName.Home,
			component: HomeView
		},
		{
			path: '/updates',
			name: RouteName.Updates,
			// route level code-splitting
			// this generates a separate chunk (About.[hash].js) for this route
			// which is lazy-loaded when the route is visited.
			component: () => import("../views/UpdatesView.vue"),
		},

		// Items
		{
			path: '/items/',
			name: RouteName.ItemList,
			component: () => import("../views/Item/ItemListView.vue"),
		},
		{
			path: '/items/patch/:patch',
			name: RouteName.ItemPatch,
			component: () => import("../views/Item/ItemPatchView.vue"),
		},
		{
			path: '/items/:itemId',
			name: RouteName.ItemHistory,
			component: () => import("../views/Item/ItemHistoryView.vue"),
		},

		// Quests
		{
			path: '/quests/',
			name: RouteName.QuestList,
			component: () => import("../views/Quest/QuestListView.vue"),
		},
		{
			path: '/quests/patch/:patch',
			name: RouteName.QuestPatch,
			component: () => import("../views/Quest/QuestPatchView.vue"),
		},
		{
			path: '/quests/:questId',
			name: RouteName.QuestHistory,
			component: () => import("../views/Quest/QuestHistoryView.vue"),
		},

		// States
		{
			path: '/states/',
			name: RouteName.StateList,
			component: () => import("../views/State/StateListView.vue"),
		},
		{
			path: '/states/patch/:patch',
			name: RouteName.StatePatch,
			component: () => import("../views/State/StatePatchView.vue"),
		},
		{
			path: '/states/:stateId',
			name: RouteName.StateHistory,
			component: () => import("../views/State/StateHistoryView.vue"),
		},
	]
})

export default router
