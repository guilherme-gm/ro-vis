<script setup lang="ts">
import DiffedValue from '@/components/DiffedValue.vue';
import MoveInfoCell from '@/components/Item/MoveInfoCell.vue';
import MoveInfoDiffCell from '@/components/Item/MoveInfoDiffCell.vue';
import ListingBase from '@/components/ListingBase.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import type { Item } from '@/models/Item';
import { useItems, type PatchItem } from '@/services/items';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const patch = route.params.patch as string;

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

const fields = ref<[string, keyof Item][]>([
	['Identified name', 'IdentifiedName'],
	['Identified description', 'IdentifiedDescription'],
	['Identified sprite', 'IdentifiedSprite'],
	['Unidentified name', 'UnidentifiedName'],
	['Unidentified description', 'UnidentifiedDescription'],
	['Unidentified sprite', 'UnidentifiedSprite'],
	['Slot count', 'SlotCount'],
	['Is book', 'IsBook'],
	['Can use buying store', 'CanUseBuyingStore'],
	['Card prefix', 'CardPrefix'],
	['Card postfix', 'CardPostfix'],
	['Card illustration', 'CardIllustration'],
	['Class num', 'ClassNum'],
]);

loadPage(1);

const newItems = computed(() => list.value.filter((v) => v.previous === null));
const updatedItems = computed(() => list.value.filter((v) => v.previous !== null && v.current !== null));
const deletedItems = computed(() => list.value.filter((v) => v.previous !== null && v.current === null));
</script>

<template>
	<ListingBase
		title="Items in patch"
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
				<table class="table table-striped table-sm">
					<tbody>
						<tr>
							<th>Info</th>
							<th>New</th>
						</tr>
						<tr v-for="(info) of fields" :key="info[1]">
							<th>{{ info[0] }}</th>
							<td><pre class="pre-preserve">{{ val.current?.[info[1]] ?? "-" }}</pre></td>
						</tr>
						<tr>
							<th>Move Info</th>
							<td><MoveInfoCell :info="val.current?.MoveInfo ?? null" /></td>
						</tr>
					</tbody>
				</table>
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
				<table class="table table-striped table-sm">
					<tbody>
						<tr>
							<th>Info</th>
							<th>Previous</th>
							<th>New</th>
							<th>Diff</th>
						</tr>
						<tr v-for="(info) of fields" :key="info[1]">
							<th>{{ info[0] }}</th>
							<td><pre class="pre-preserve">{{ val.previous?.[info[1]] ?? "-" }}</pre></td>
							<td><pre class="pre-preserve">{{ val.current?.[info[1]] ?? "-" }}</pre></td>
							<td><DiffedValue :from="val.previous![info[1]]" :to="val.current![info[1]]" /></td>
						</tr>
						<tr>
							<th>Move Info</th>
							<td><MoveInfoCell :info="val.previous?.MoveInfo ?? null" /></td>
							<td><MoveInfoCell :info="val.current?.MoveInfo ?? null" /></td>
							<td><MoveInfoDiffCell :from="val.previous?.MoveInfo ?? null" :to="val.current?.MoveInfo ?? null" /></td>
						</tr>
					</tbody>
				</table>
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
				<table class="table table-striped table-sm">
					<tbody>
						<tr>
							<th>Info</th>
							<th>Previous</th>
						</tr>
						<tr v-for="(info) of fields" :key="info[1]">
							<th>{{ info[0] }}</th>
							<td><pre class="pre-preserve">{{ val.previous?.[info[1]] ?? "-" }}</pre></td>
						</tr>
						<tr>
							<th>Move Info</th>
							<td><MoveInfoCell :info="val.previous?.MoveInfo ?? null" /></td>
						</tr>
					</tbody>
				</table>
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted items in this page</p>
	</ListingBase>
</template>
