<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import type { MinMap } from '@/models/Map';
import { RouteName } from '@/router/RouteName';
import { MapApi } from '@/services/MapApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { ref } from 'vue';

document.title = `Map List - ROVis`;

const {
	state,
	total,
	getItems,
} = MapApi.use();

const list = ref<MinMap[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItems(page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		title="Maps"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsListGroup :flush="true">
			<BsListGroupItem
				v-for="(map) in list"
				:key="`${map.MapID}-${map.LastUpdate}`"
			>
				<BsLink
					:to="{ name: RouteName.MapHistory, params: { mapId: map.MapID } }"
					target="_blank"
				>
					#{{ map.MapID }} - {{ map.Name }} (Last updated: {{ map.LastUpdate }})
					<BIconBoxArrowUpRight />
				</BsLink>
			</BsListGroupItem>
		</BsListGroup>
	</ListingBase>
</template>
