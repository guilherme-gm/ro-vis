<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import StateCompareTable from '@/components/State/StateCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import { RouteName } from '@/router/RouteName';
import { StateApi, type StatePatch } from '@/services/StateApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const patch = route.params.patch as string;

document.title = `Patch #${patch} states - ROVis`;

const {
	state,
	total,
	getPatchItems,
} = StateApi.use();

const list = ref<StatePatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getPatchItems(patch, page);
}

loadPage(1);

const newItems = computed(() => list.value.filter((v) => v.From === null));
const updatedItems = computed(() => list.value.filter((v) => v.From !== null && v.To !== null));
const deletedItems = computed(() => list.value.filter((v) => v.From !== null && v.To === null));
</script>

<template>
	<ListingBase
		:title="`States in patch ${patch}`"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<h4>New states</h4>
		<BsAccordion v-if="newItems.length > 0">
			<BsAccordionItem
				v-for="(val) in newItems"
				:key="val.To?.Data.Id"
				:title="`#${val.To?.Data.Id} - ${val?.To?.Data.Constant}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.StateHistory, params: { stateId: val.To?.Data.Id ?? val.From?.Data.Id } }"
						target="_blank"
					>
						View State history <BIconBoxArrowUpRight />
					</BsLink>
				</p>
				<StateCompareTable :current="val.To?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No new states in this page</p>

		<h4 class="mt-3">Updated States</h4>
		<BsAccordion v-if="updatedItems.length > 0">
			<BsAccordionItem
				v-for="(val) in updatedItems"
				:key="val.To?.Data.Id ?? val.From?.Data.Id"
				:title="`#${val.To?.Data.Id ?? val.From?.Data.Id} - ${val?.To?.Data.Constant ?? val?.From?.Data.Constant}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.StateHistory, params: { stateId: val.To?.Data.Id ?? val.From?.Data.Id } }"
						target="_blank"
					>
						View State history <BIconBoxArrowUpRight />
					</BsLink>
				</p>
				<StateCompareTable :current="val.To?.Data" :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No updated states in this page</p>

		<h4 class="mt-3">Deleted States</h4>
		<BsAccordion v-if="deletedItems.length > 0">
			<BsAccordionItem
				v-for="(val) in deletedItems"
				:key="val.From?.Data.Id"
				:title="`#${val.From?.Data.Id} - ${val?.From?.Data.Constant}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.StateHistory, params: { stateId: val.To?.Data.Id ?? val.From?.Data.Id } }"
						target="_blank"
					>
						View State history <BIconBoxArrowUpRight />
					</BsLink>
				</p>
				<StateCompareTable :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted states in this page</p>
	</ListingBase>
</template>
