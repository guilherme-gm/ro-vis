<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import type { MinI18n } from '@/models/I18n';
import { RouteName } from '@/router/RouteName';
import { I18nApi } from '@/services/I18nApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { ref } from 'vue';

document.title = `i18n List - ROVis`;

const {
	state,
	total,
	getItems,
} = I18nApi.use();

const list = ref<MinI18n[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItems(page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		title="i18n Entries"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsListGroup :flush="true">
			<BsListGroupItem
				v-for="(entry) in list"
				:key="`${entry.I18nId}`"
			>
				<BsLink
					:to="{ name: RouteName.I18nHistory, params: { i18nId: entry.I18nId } }"
					target="_blank"
				>
					#{{ entry.I18nId }} - {{ entry.PtBrText }} (Last updated: {{ entry.LastUpdate }})
					<BIconBoxArrowUpRight />
				</BsLink>
			</BsListGroupItem>
		</BsListGroup>
	</ListingBase>
</template>
