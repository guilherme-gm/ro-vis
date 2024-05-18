<script setup lang="ts">
import BsNavbar from "@/components/bootstrap/BsNavbar.vue";
import BsNavLink from "./bootstrap/BsNavLink.vue";
import { useRoute } from "vue-router";
import { computed } from "vue";
import { RouteName } from "@/router/RouteName";

const route = useRoute();
const currentRoute = computed(() => route.name);
</script>

<template>
	<header>
		<BsNavbar brand-text="RO Vis">
			<template #nav-items>
				<BsNavLink :to="{ name: RouteName.Home }" label="Home" :current-page="currentRoute === RouteName.Home" />
			</template>
		</BsNavbar>
	</header>
</template>
