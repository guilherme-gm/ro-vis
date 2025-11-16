<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import SkillCompareTable from '@/components/Skill/SkillCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import { SkillApi, type SkillPatch } from '@/services/SkillApi';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const skillId = route.params.skillId as string;

document.title = `Skill #${skillId} history - ROVis`;

const {
	state,
	historyTotal,
	getItemHistory,
} = SkillApi.use();

const list = ref<SkillPatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItemHistory(skillId, page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		:title="`History for Skill ID #${skillId}`"
		:total="historyTotal"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsAccordion v-if="list.length > 0">
			<BsAccordionItem
				v-for="(val) in list"
				:key="val.To?.Data.SkillID ?? val.From?.Data.SkillID"
				:title="`Patch ${val.To?.Update ?? '[unknown]'}`"
				:expanded="true"
			>
				<SkillCompareTable :current="val.To?.Data" :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<span v-else>Skill not found.</span>
	</ListingBase>
</template>
