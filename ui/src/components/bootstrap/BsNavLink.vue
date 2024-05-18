<script setup lang="ts">
import { RouterLink } from "vue-router";

withDefaults(defineProps<{
	label: string;
	to?: string;
	href?: string;
	currentPage?: boolean;
}>(), {
	label: '',
	currentPage: false,
});
</script>

<template>
	<li class="nav-item">
		<RouterLink
			v-if="to"
			:to="to"
			:aria-current="currentPage ? 'page' : 'false'"
			class="nav-link"
			:class="{ active: currentPage }"
		>
			{{ label }}
		</RouterLink>
		<a
			v-else
			:href="href"
			:aria-current="currentPage ? 'page' : 'false'"
			class="nav-link"
			:class="{ active: currentPage }"
		>
			{{ label }}
		</a>
	</li>
</template>
