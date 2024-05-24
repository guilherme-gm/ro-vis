<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import type { Update } from '@/models/Update';
import { useUpdates } from '@/services/updates';
import { ref } from 'vue';

const {
	state,
	total,
	getUpdates,
} = useUpdates();

const list = ref<Update[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getUpdates(page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		title="Update history"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsAccordion>
			<BsAccordionItem
				v-for="(val) in list"
				:key="val.id"
				:title="val.id"
			>
				<h5>Changed files:</h5>
				<BsListGroup :flush="true">
					<BsListGroupItem
						v-for="(update) in val.updates"
						:key="`${val.id}-${update.file}`"
					>
						{{ update.file }}
					</BsListGroupItem>
				</BsListGroup>
			</BsAccordionItem>
		</BsAccordion>
	</ListingBase>
</template>
