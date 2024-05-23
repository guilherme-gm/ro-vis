<script setup lang="ts">
defineProps<{
	pages: number[];
	currentPage: number;
	description: string;
	center?: boolean;
}>();

const emit = defineEmits(['changed']);
</script>

<template>
	<nav :aria-label="description">
		<ul class="pagination" :class="{ 'justify-content-center': center }">
			<li
				v-for="(page) in pages"
				:key="`page-${page}`"
				class="page-item"
				:class="{ active: (page === currentPage) }"
			>
				<button class="page-link" @click="emit('changed', page)">{{ page }}</button>
			</li>
		</ul>
	</nav>
</template>
