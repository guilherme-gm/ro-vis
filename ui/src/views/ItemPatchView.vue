<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import type { Item } from '@/models/Item';
import { useItems } from '@/services/items';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const patch = route.params.patch as string;

const {
	state,
	total,
	getPatchItems,
} = useItems();

const list = ref<Item[]>([]);
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
</script>

<template>
	<ListingBase
		title="Items in patch"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsAccordion>
			<BsAccordionItem
				v-for="(val) in list"
				:key="val.ItemId"
				:title="`#${val.ItemId} - ${val.IdentifiedName} (${val.MoveInfo?.commentName ?? ''})`"
			>
				<table class="table table-striped table-sm">
					<tbody>
						<tr>
							<th>Info</th>
							<th>Previous</th>
							<th>New</th>
						</tr>
						<tr v-for="(info) of fields" :key="info[1]">
							<th>{{ info[0] }}</th>
							<td><pre>{{ (Array.isArray(val[info[1]]) ? (val[info[1]] as string[]).join('\n') : val[info[1]]) }}</pre></td>
							<td><pre>{{ (Array.isArray(val[info[1]]) ? (val[info[1]] as string[]).join('\n') : val[info[1]]) }}</pre></td>
						</tr>
					</tbody>
				</table>
			</BsAccordionItem>
		</BsAccordion>
	</ListingBase>
</template>
