<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import QuestCompareTable from '@/components/Quest/QuestCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import { QuestApi, type QuestPatch } from '@/services/QuestApi';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const questId = route.params.questId as string;

document.title = `Item #${questId} history - ROVis`;

const {
	state,
	historyTotal,
	getItemHistory,
} = QuestApi.use();

const list = ref<QuestPatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItemHistory(questId, page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		:title="`History for Quest ID #${questId}`"
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
				<QuestCompareTable :current="val.To" :previous="val.From" />
			</BsAccordionItem>
		</BsAccordion>
		<span v-else>Quest not found.</span>
	</ListingBase>
</template>
