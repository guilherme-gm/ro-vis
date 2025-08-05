<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import MapCompareTable from '@/components/Map/MapCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import { MapApi, type MapPatch } from '@/services/MapApi';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const mapId = route.params.mapId as string;

document.title = `Map #${mapId} history - ROVis`;

const {
	state,
	historyTotal,
	getItemHistory,
} = MapApi.use();

const list = ref<MapPatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItemHistory(mapId, page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		:title="`History for Map ID #${mapId}`"
		:total="historyTotal"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsAccordion v-if="list.length > 0">
			<BsAccordionItem
				v-for="(val) in list"
				:key="val.To?.Data.Id ?? val.From?.Data.Id"
				:title="`Patch ${val.To?.Update ?? '[unknown]'}`"
				:expanded="true"
			>
				<MapCompareTable :current="val.To?.Data" :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<span v-else>Map not found.</span>
	</ListingBase>
</template>
