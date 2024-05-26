<script setup lang="ts">
import ItemCompareTable from '@/components/Item/ItemCompareTable.vue';
import ListingBase from '@/components/ListingBase.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import { RouteName } from '@/router/RouteName';
import { useItems, type PatchItem } from '@/services/items';
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
} = useItems();

const list = ref<PatchItem[]>([]);
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
				:key="val.current?.ItemId"
				:title="`#${val.current?.ItemId} - ${val?.current?.IdentifiedName} (${val?.current?.MoveInfo?.commentName ?? ''})`"
			>
				<p>
					<BsLink
						:to="{ name: RouteName.ItemHistory, params: { itemId: val.current?.ItemId ?? val.previous?.ItemId } }"
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
				:key="val.current?.ItemId ?? val.previous?.ItemId"
				:title="`#${val.current?.ItemId ?? val.previous?.ItemId} - ${val?.current?.IdentifiedName ?? val?.previous?.IdentifiedName} (${val?.current?.MoveInfo?.commentName ?? ''})`"
			>
				<ItemCompareTable :current="val.current" :previous="val.previous" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No updated items in this page</p>

		<h4 class="mt-3">Deleted Items</h4>
		<BsAccordion v-if="deletedItems.length > 0">
			<BsAccordionItem
				v-for="(val) in deletedItems"
				:key="val.previous?.ItemId"
				:title="`#${val.previous?.ItemId} - ${val?.previous?.IdentifiedName} (${val?.previous?.MoveInfo?.commentName ?? ''})`"
			>
				<ItemCompareTable :previous="val.previous" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted items in this page</p>
	</ListingBase>
</template>
