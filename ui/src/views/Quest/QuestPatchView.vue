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

const newQuests = computed(() => list.value.filter((v) => v.From === null));
const updatedQuests = computed(() => list.value.filter((v) => v.From !== null && v.To !== null));
const deletedQuests = computed(() => list.value.filter((v) => v.From !== null && v.To === null));
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
				:key="val.To?.Id"
				:title="`#${val.To?.Id} - ${val?.To?.Title}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.QuestHistory, params: { questId: val.To?.Id ?? val.From?.Id } }"
						target="_blank"
					>
						View Quest history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<QuestCompareTable :current="val.To" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No new quests in this page</p>

		<h4 class="mt-3">Updated Quests</h4>
		<BsAccordion v-if="updatedQuests.length > 0">
			<BsAccordionItem
				v-for="(val) in updatedQuests"
				:key="val.To?.Id ?? val.From?.Id"
				:title="`#${val.To?.Id ?? val.From?.Id} - ${val?.To?.Title ?? val?.From?.Title}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.QuestHistory, params: { questId: val.To?.Id ?? val.From?.Id } }"
						target="_blank"
					>
						View Quest history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<QuestCompareTable :current="val.To" :previous="val.From" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No updated quests in this page</p>

		<h4 class="mt-3">Deleted Quests</h4>
		<BsAccordion v-if="deletedQuests.length > 0">
			<BsAccordionItem
				v-for="(val) in deletedQuests"
				:key="val.From?.Id"
				:title="`#${val.From?.Id} - ${val?.From?.Title}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.lastUpdate }}
					<BsLink
						:to="{ name: RouteName.QuestHistory, params: { questId: val.To?.Id ?? val.From?.Id } }"
						target="_blank"
					>
						View Quest history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<QuestCompareTable :previous="val.From" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted quests in this page</p>
	</ListingBase>
</template>
