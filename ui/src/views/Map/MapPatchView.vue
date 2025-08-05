<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import MapCompareTable from '@/components/Map/MapCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import { RouteName } from '@/router/RouteName';
import { MapApi, type MapPatch } from '@/services/MapApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const patch = route.params.patch as string;

document.title = `Patch #${patch} maps - ROVis`;

const {
	state,
	total,
	getPatchItems,
} = MapApi.use();

const list = ref<MapPatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getPatchItems(patch, page);
}

loadPage(1);

const newMaps = computed(() => list.value.filter((v) => !v.From));
const updatedMaps = computed(() => list.value.filter((v) => !!v.From && !!v.To));
const deletedMaps = computed(() => list.value.filter((v) => !!v.From && !v.To));
</script>

<template>
	<ListingBase
		:title="`Maps in patch ${patch}`"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<h4>New maps</h4>
		<BsAccordion v-if="newMaps.length > 0">
			<BsAccordionItem
				v-for="(val) in newMaps"
				:key="val.To?.Data.Id"
				:title="`#${val.To?.Data.Id} - ${val?.To?.Data.Name}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.MapHistory, params: { mapId: val.To?.Data.Id ?? val.From?.Data.Id } }"
						target="_blank"
					>
						View Map history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<MapCompareTable :current="val.To?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No new maps in this page</p>

		<h4 class="mt-3">Updated maps</h4>
		<BsAccordion v-if="updatedMaps.length > 0">
			<BsAccordionItem
				v-for="(val) in updatedMaps"
				:key="val.To?.Data.Id ?? val.From?.Data.Id"
				:title="`#${val.To?.Data.Id ?? val.From?.Data.Id} - ${val?.To?.Data.Name ?? val?.From?.Data.Name}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.MapHistory, params: { mapId: val.To?.Data.Id ?? val.From?.Data.Id } }"
						target="_blank"
					>
						View Map history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<MapCompareTable :current="val.To?.Data" :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No updated maps in this page</p>

		<h4 class="mt-3">Deleted maps</h4>
		<BsAccordion v-if="deletedMaps.length > 0">
			<BsAccordionItem
				v-for="(val) in deletedMaps"
				:key="val.From?.Data.Id"
				:title="`#${val.From?.Data.Id} - ${val?.From?.Data.Name}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.MapHistory, params: { mapId: val.To?.Data.Id ?? val.From?.Data.Id } }"
						target="_blank"
					>
						View Map history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<MapCompareTable :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted maps in this page</p>
	</ListingBase>
</template>
