<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import I18nCompareTable from '@/components/I18n/I18nCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import { I18nApi, type I18nPatch } from '@/services/I18nApi';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const i18nId = route.params.i18nId as string;

document.title = `i18n #${i18nId} history - ROVis`;

const {
	state,
	historyTotal,
	getItemHistory,
} = I18nApi.use();

const list = ref<I18nPatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItemHistory(i18nId, page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		:title="`History for i18n ID #${i18nId}`"
		:total="historyTotal"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsAccordion v-if="list.length > 0">
			<BsAccordionItem
				v-for="(val) in list"
				:key="val.To?.Data.I18nId ?? val.From?.Data.I18nId"
				:title="`Update ${val.To?.Update ?? val.From?.Update ?? 'unknown'}`"
				:expanded="true"
			>
				<I18nCompareTable :current="val.To?.Data" :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<span v-else>i18n entry not found.</span>
	</ListingBase>
</template>
