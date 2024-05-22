<script setup lang="ts">
import { RouterView } from "vue-router";
import SiteHeader from "./components/SiteHeader.vue";
import { ref } from "vue";
import BsBorderSpinner from "./components/bootstrap/Spinner/BsBorderSpinner.vue";
import { BIconXLg } from "bootstrap-icons-vue";

const loading = ref(false);
const errorMesage = ref('');
</script>

<template>
	<SiteHeader />

	<main v-if="loading" class="d-flex flex-fill">
		<div v-if="!errorMesage" class="m-auto d-flex flex-column">
			<BsBorderSpinner v-if="loading" class="m-auto" />
			<span class="pt-3">Loading RO Vis...</span>
		</div>
		<div v-else class="m-auto d-flex flex-column">
			<BIconXLg class="fs-2 m-auto text-danger" />
			<p class="pt-3">There was an error loading RO Vis...</p>
			<p>Try refreshing the page or try again later.</p>
			<p><strong>Details:</strong> {{ errorMesage }}</p>
		</div>
	</main>
	<main v-else class="flex-shrink-0 h-100 p-4">
		<RouterView />
	</main>
</template>

<style scoped></style>
