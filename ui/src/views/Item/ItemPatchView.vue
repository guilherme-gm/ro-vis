<script setup lang="ts">
import ItemCompareTable from '@/components/Item/ItemCompareTable.vue';
import ListingBase from '@/components/ListingBase.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import { RouteName } from '@/router/RouteName';
import { ItemApi, type ItemPatch } from '@/services/ItemApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const patch = route.params.patch as string;

document.title = `Patch #${patch} items - ROVis`;

const {
	state,
	total,
	getPatchItems,
} = ItemApi.use();

const list = ref<ItemPatch[]>([]);
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
		:title="`Items in patch ${patch}`"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<h4>New items</h4>
		<BsAccordion v-if="newItems.length > 0">
			<BsAccordionItem
				v-for="(val) in newItems"
				:key="val.current?.Id"
				:title="`#${val.current?.Id} - ${val?.current?.IdentifiedName} (${val?.current?.MoveInfo?.commentName ?? ''})`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.ItemHistory, params: { itemId: val.current?.Id ?? val.previous?.Id } }"
						target="_blank"
					>
						View Item history <BIconBoxArrowUpRight />
					</BsLink>
				</p>
				<ItemCompareTable :current="val.current" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No new items in this page</p>

		<h4 class="mt-3">Updated Items</h4>
		<BsAccordion v-if="updatedItems.length > 0">
			<BsAccordionItem
				v-for="(val) in updatedItems"
				:key="val.current?.Id ?? val.previous?.Id"
				:title="`#${val.current?.Id ?? val.previous?.Id} - ${val?.current?.IdentifiedName ?? val?.previous?.IdentifiedName} (${val?.current?.MoveInfo?.commentName ?? ''})`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.ItemHistory, params: { itemId: val.current?.Id ?? val.previous?.Id } }"
						target="_blank"
					>
						View Item history <BIconBoxArrowUpRight />
					</BsLink>
				</p>
				<ItemCompareTable :current="val.current" :previous="val.previous" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No updated items in this page</p>

		<h4 class="mt-3">Deleted Items</h4>
		<BsAccordion v-if="deletedItems.length > 0">
			<BsAccordionItem
				v-for="(val) in deletedItems"
				:key="val.previous?.Id"
				:title="`#${val.previous?.Id} - ${val?.previous?.IdentifiedName} (${val?.previous?.MoveInfo?.commentName ?? ''})`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.ItemHistory, params: { itemId: val.current?.Id ?? val.previous?.Id } }"
						target="_blank"
					>
						View Item history <BIconBoxArrowUpRight />
					</BsLink>
				</p>
				<ItemCompareTable :previous="val.previous" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted items in this page</p>
	</ListingBase>
</template>
