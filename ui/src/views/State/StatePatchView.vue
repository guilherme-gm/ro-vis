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

const newItems = computed(() => list.value.filter((v) => v.previous === null));
const updatedItems = computed(() => list.value.filter((v) => v.previous !== null && v.current !== null));
const deletedItems = computed(() => list.value.filter((v) => v.previous !== null && v.current === null));
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
				:key="val.current?.Id"
				:title="`#${val.current?.Id} - ${val?.current?.Constant}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.StateHistory, params: { stateId: val.current?.Id ?? val.previous?.Id } }"
						target="_blank"
					>
						View State history <BIconBoxArrowUpRight />
					</BsLink>
				</p>
				<StateCompareTable :current="val.current" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No new states in this page</p>

		<h4 class="mt-3">Updated States</h4>
		<BsAccordion v-if="updatedItems.length > 0">
			<BsAccordionItem
				v-for="(val) in updatedItems"
				:key="val.current?.Id ?? val.previous?.Id"
				:title="`#${val.current?.Id ?? val.previous?.Id} - ${val?.current?.Constant ?? val?.previous?.Constant}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.StateHistory, params: { stateId: val.current?.Id ?? val.previous?.Id } }"
						target="_blank"
					>
						View State history <BIconBoxArrowUpRight />
					</BsLink>
				</p>
				<StateCompareTable :current="val.current" :previous="val.previous" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No updated states in this page</p>

		<h4 class="mt-3">Deleted States</h4>
		<BsAccordion v-if="deletedItems.length > 0">
			<BsAccordionItem
				v-for="(val) in deletedItems"
				:key="val.previous?.Id"
				:title="`#${val.previous?.Id} - ${val?.previous?.Constant}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.StateHistory, params: { stateId: val.current?.Id ?? val.previous?.Id } }"
						target="_blank"
					>
						View State history <BIconBoxArrowUpRight />
					</BsLink>
				</p>
				<StateCompareTable :previous="val.previous" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted states in this page</p>
	</ListingBase>
</template>
