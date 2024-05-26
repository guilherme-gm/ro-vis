<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import type { MinQuest } from '@/models/Quest';
import { RouteName } from '@/router/RouteName';
import { QuestApi } from '@/services/QuestApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { ref } from 'vue';

document.title = `Quest List - ROVis`;

const {
	state,
	total,
	getItems,
} = QuestApi.use();

const list = ref<MinQuest[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItems(page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		title="Quests"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsListGroup :flush="true">
			<BsListGroupItem
				v-for="(quest) in list"
				:key="`${quest.id}-${quest.patch}`"
			>
				<BsLink
					:to="{ name: RouteName.QuestHistory, params: { questId: quest.id } }"
					target="_blank"
				>
					#{{ quest.id }} - {{ quest.Title }} (Last updated: {{ quest.patch }})
					<BIconBoxArrowUpRight />
				</BsLink>
			</BsListGroupItem>
		</BsListGroup>
	</ListingBase>
</template>
