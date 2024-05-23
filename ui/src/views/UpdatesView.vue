<script setup lang="ts">
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import { useUpdates } from '@/services/updates';
import { LoadState } from "@/services/LoadState";
import type { Update } from '@/models/Update';
import { computed, ref } from 'vue';
import BsBorderSpinner from '@/components/bootstrap/Spinner/BsBorderSpinner.vue';
import { BIconXLg } from 'bootstrap-icons-vue';
import BsPagination from '@/components/bootstrap/Pagination/BsPagination.vue';

const {
	state,
	total,
	getUpdates,
} = useUpdates();

const list = ref<Update[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getUpdates(page);
}

const pages = computed(() => {
	const nums = [];

	let pageNum = 1;
	for (let i = 0; i <= total.value; i += 100) {
		nums.push(pageNum);
		pageNum++;
	}

	return nums;
});

loadPage(1);
</script>

<template>
	<h2>Updates History</h2>
	<p v-if="total >= 0">
		<strong>{{ total }}</strong> items to show.
	</p>

	<div v-if="state === LoadState.Loaded">
		<BsAccordion>
			<BsAccordionItem
				v-for="(val) in list"
				:key="val.id"
				:title="val.id"
			>
				<h5>Changed files:</h5>
				<BsListGroup :flush="true">
					<BsListGroupItem
						v-for="(update) in val.updates"
						:key="`${val.id}-${update.file}`"
					>
						{{ update.file }}
					</BsListGroupItem>
				</BsListGroup>
			</BsAccordionItem>
		</BsAccordion>
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
		description="Update list page"
		:pages="pages"
		:current-page="currentPage"
		:center="true"
		@changed="loadPage"
		class="mt-3"
	/>
</template>
