<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import QuestCompareTable from '@/components/Quest/QuestCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import { RouteName } from '@/router/RouteName';
import { QuestApi, type QuestPatch } from '@/services/QuestApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const patch = route.params.patch as string;

document.title = `Patch #${patch} quests - ROVis`;

const {
	state,
	total,
	getPatchItems,
} = QuestApi.use();

const list = ref<QuestPatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getPatchItems(patch, page);
}

loadPage(1);

const newQuests = computed(() => list.value.filter((v) => v.previous === null));
const updatedQuests = computed(() => list.value.filter((v) => v.previous !== null && v.current !== null));
const deletedQuests = computed(() => list.value.filter((v) => v.previous !== null && v.current === null));
</script>

<template>
	<ListingBase
		:title="`Quests in patch ${patch}`"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<h4>New quests</h4>
		<BsAccordion v-if="newQuests.length > 0">
			<BsAccordionItem
				v-for="(val) in newQuests"
				:key="val.current?.Id"
				:title="`#${val.current?.Id} - ${val?.current?.Title}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.QuestHistory, params: { questId: val.current?.Id ?? val.previous?.Id } }"
						target="_blank"
					>
						View Quest history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<QuestCompareTable :current="val.current" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No new quests in this page</p>

		<h4 class="mt-3">Updated Quests</h4>
		<BsAccordion v-if="updatedQuests.length > 0">
			<BsAccordionItem
				v-for="(val) in updatedQuests"
				:key="val.current?.Id ?? val.previous?.Id"
				:title="`#${val.current?.Id ?? val.previous?.Id} - ${val?.current?.Title ?? val?.previous?.Title}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.QuestHistory, params: { questId: val.current?.Id ?? val.previous?.Id } }"
						target="_blank"
					>
						View Quest history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<QuestCompareTable :current="val.current" :previous="val.previous" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No updated quests in this page</p>

		<h4 class="mt-3">Deleted Quests</h4>
		<BsAccordion v-if="deletedQuests.length > 0">
			<BsAccordionItem
				v-for="(val) in deletedQuests"
				:key="val.previous?.Id"
				:title="`#${val.previous?.Id} - ${val?.previous?.Title}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.QuestHistory, params: { questId: val.current?.Id ?? val.previous?.Id } }"
						target="_blank"
					>
						View Quest history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<QuestCompareTable :previous="val.previous" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted quests in this page</p>
	</ListingBase>
</template>
