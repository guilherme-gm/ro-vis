<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import StateCompareTable from '@/components/State/StateCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import { StateApi, type StatePatch } from '@/services/StateApi';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const stateId = route.params.stateId as string;

document.title = `State #${stateId} history - ROVis`;

const {
	state,
	historyTotal,
	getItemHistory,
} = StateApi.use();

const list = ref<StatePatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItemHistory(stateId, page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		:title="`History for state ID #${stateId}`"
		:total="historyTotal"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsAccordion v-if="list.length > 0">
			<BsAccordionItem
				v-for="(val) in list"
				:key="val.To?.Id ?? val.From?.Id"
				:title="`Patch ${val.To?.Patch ?? '[unknown]'}`"
				:expanded="true"
			>
				<StateCompareTable :current="val.To" :previous="val.From" />
			</BsAccordionItem>
		</BsAccordion>
		<span v-else>State not found.</span>
	</ListingBase>
</template>
