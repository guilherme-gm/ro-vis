<script setup lang="ts">
import ItemCompareTable from '@/components/Item/ItemCompareTable.vue';
import ListingBase from '@/components/ListingBase.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import { ItemApi, type ItemPatch } from '@/services/ItemApi';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const itemId = route.params.itemId as string;

document.title = `Item #${itemId} history - ROVis`;

const {
	state,
	historyTotal,
	getItemHistory,
} = ItemApi.use();

const list = ref<ItemPatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItemHistory(itemId, page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		:title="`History for item ID #${itemId}`"
		:total="historyTotal"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsAccordion v-if="list.length > 0">
			<BsAccordionItem
				v-for="(val) in list"
				:key="val.current?.Id ?? val.previous?.Id"
				:title="`Patch ${val.current?.Patch ?? '[unknown]'}`"
				:expanded="true"
			>
				<ItemCompareTable :current="val.current" :previous="val.previous" />
			</BsAccordionItem>
		</BsAccordion>
		<span v-else>Item not found.</span>
	</ListingBase>
</template>
