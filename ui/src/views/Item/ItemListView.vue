<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import type { MinItem } from '@/models/Item';
import { RouteName } from '@/router/RouteName';
import { ItemApi } from '@/services/ItemApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { ref } from 'vue';

document.title = `Item list - ROVis`;

const {
	state,
	total,
	getItems,
} = ItemApi.use();

const list = ref<MinItem[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItems(page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		title="Items"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsListGroup :flush="true">
			<BsListGroupItem
				v-for="(item) in list"
				:key="`${item.Id}-${item.Patch}`"
			>
				<BsLink
					:to="{ name: RouteName.ItemHistory, params: { itemId: item.Id } }"
					target="_blank"
				>
					#{{ item.Id }} - {{ item.IdentifiedName }} (Last updated: {{ item.Patch }})
					<BIconBoxArrowUpRight />
				</BsLink>
			</BsListGroupItem>
		</BsListGroup>
	</ListingBase>
</template>
