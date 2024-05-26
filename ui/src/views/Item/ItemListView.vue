<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import type { MinItem } from '@/models/Item';
import { RouteName } from '@/router/RouteName';
import { useItems } from '@/services/items';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const itemId = route.params.itemId as string;

document.title = `Item #${itemId} history - ROVis`;

const {
	state,
	total,
	getItems,
} = useItems();

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
				:key="`${item.id}-${item.patch}`"
			>
				<BsLink
					:to="{ name: RouteName.ItemHistory, params: { itemId: item.id } }"
					target="_blank"
				>
					#{{ item.id }} - {{ item.IdentifiedName }} (Last updated: {{ item.patch }})
					<BIconBoxArrowUpRight />
				</BsLink>
			</BsListGroupItem>
		</BsListGroup>
	</ListingBase>
</template>
