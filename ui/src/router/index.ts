import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { RouteName } from './RouteName'

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
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
		{
			path: '/items/patch/:patch',
			name: RouteName.ItemPatch,
			component: () => import("../views/ItemPatchView.vue"),
		},
		{
			path: '/items/:itemId',
			name: RouteName.ItemHistory,
			component: () => import("../views/ItemHistoryView.vue"),
		},
	]
})

export default router
