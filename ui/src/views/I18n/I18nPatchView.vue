<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import I18nCompareTable from '@/components/I18n/I18nCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import { RouteName } from '@/router/RouteName';
import { I18nApi, type I18nPatch } from '@/services/I18nApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const patch = route.params.patch as string;

document.title = `Patch #${patch} i18n entries - ROVis`;

const {
	state,
	total,
	getPatchItems,
} = I18nApi.use();

const list = ref<I18nPatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getPatchItems(patch, page);
}

loadPage(1);

const newEntries = computed(() => list.value.filter((v) => !v.From));
const updatedEntries = computed(() => list.value.filter((v) => !!v.From && !!v.To));
const deletedEntries = computed(() => list.value.filter((v) => !!v.From && !v.To));
</script>

<template>
	<ListingBase
		:title="`i18n entries in patch ${patch}`"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<h4>New entries</h4>
		<BsAccordion v-if="newEntries.length > 0">
			<BsAccordionItem
				v-for="(val) in newEntries"
				:key="val.To?.Data.I18nId"
				:title="`#${val.To?.Data.I18nId} - ${val.To?.Data.PtBrText} [${val.To?.Data.ContainerFile}]`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.I18nHistory, params: { i18nId: val.To?.Data.I18nId } }"
						target="_blank"
					>
						View i18n history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<I18nCompareTable :current="val.To?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No new i18n entries in this page</p>

		<h4 class="mt-3">Updated Entries</h4>
		<BsAccordion v-if="updatedEntries.length > 0">
			<BsAccordionItem
				v-for="(val) in updatedEntries"
				:key="val.To?.Data.I18nId ?? val.From?.Data.I18nId"
				:title="`#${val.To?.Data.I18nId ?? val.From?.Data.I18nId} - ${val.To?.Data.PtBrText ?? val.From?.Data.PtBrText} [${val.To?.Data.ContainerFile ?? val.From?.Data.ContainerFile}]`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.I18nHistory, params: { i18nId: val.To?.Data.I18nId ?? val.From?.Data.I18nId } }"
						target="_blank"
					>
						View i18n history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<I18nCompareTable :current="val.To?.Data" :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No updated i18n entries in this page</p>

		<h4 class="mt-3">Deleted Entries</h4>
		<BsAccordion v-if="deletedEntries.length > 0">
			<BsAccordionItem
				v-for="(val) in deletedEntries"
				:key="val.From?.Data.I18nId"
				:title="`#${val.From?.Data.I18nId} - ${val.From?.Data.PtBrText} [${val.From?.Data.ContainerFile}]`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.I18nHistory, params: { i18nId: val.From?.Data.I18nId } }"
						target="_blank"
					>
						View i18n history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<I18nCompareTable :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted i18n entries in this page</p>
	</ListingBase>
</template>
