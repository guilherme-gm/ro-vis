<script setup lang="ts">
import BsPagination from '@/components/bootstrap/Pagination/BsPagination.vue';
import BsBorderSpinner from '@/components/bootstrap/Spinner/BsBorderSpinner.vue';
import { LoadState } from "@/services/LoadState";
import { BIconXLg } from 'bootstrap-icons-vue';
import { computed } from 'vue';

const props = defineProps<{
	title: string;
	total: number;
	state: LoadState;
	currentPage: number;
}>();

const pages = computed(() => {
	const nums = [];

	let pageNum = 1;
	for (let i = 0; i <= props.total; i += 100) {
		nums.push(pageNum);
		pageNum++;
	}

	return nums;
});

const emit = defineEmits(['page-changed']);
</script>

<template>
	<h2>{{ title }}</h2>
	<p v-if="total >= 0">
		<strong>{{ total }}</strong> items to show.
	</p>

	<BsPagination
		description="list page"
		:pages="pages"
		:current-page="currentPage"
		:center="true"
		@changed="pageNum => emit('page-changed', pageNum)"
		class="mt-3"
	/>

	<div v-if="state === LoadState.Loaded">
		<slot></slot>
	</div>
	<div v-else-if="state === LoadState.Loading" class="d-flex justify-content-center">
		<BsBorderSpinner message="Loading..." />
	</div>
	<div v-else-if="state === LoadState.Error" class="m-auto d-flex flex-column">
		<BIconXLg class="fs-2 m-auto text-danger" />
		<p class="pt-3">There was an error loading this page...</p>
		<p>Try refreshing the page or try again later.</p>
	</div>

	<BsPagination
		description="list page"
		:pages="pages"
		:current-page="currentPage"
		:center="true"
		@changed="pageNum => emit('page-changed', pageNum)"
		class="mt-3"
	/>
</template>
