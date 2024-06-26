<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import type { MinState } from '@/models/State';
import { RouteName } from '@/router/RouteName';
import { StateApi } from '@/services/StateApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { ref } from 'vue';

document.title = `State list - ROVis`;

const {
	state,
	total,
	getItems,
} = StateApi.use();

const list = ref<MinState[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItems(page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		title="States"
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
					:to="{ name: RouteName.StateHistory, params: { stateId: item.Id } }"
					target="_blank"
				>
					#{{ item.Id }} - {{ item.Constant }} (Last updated: {{ item.Patch }})
					<BIconBoxArrowUpRight />
				</BsLink>
			</BsListGroupItem>
		</BsListGroup>
	</ListingBase>
</template>
