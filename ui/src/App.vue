<script setup lang="ts">
import { RouterView } from "vue-router";
import SiteHeader from "./components/SiteHeader.vue";
</script>

<template>
	<SiteHeader />

	<main class="flex-shrink-0">
		<RouterView />
	</main>
</template>

<style scoped></style>
