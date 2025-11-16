<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import SkillCompareTable from '@/components/Skill/SkillCompareTable.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import { RouteName } from '@/router/RouteName';
import { SkillApi, type SkillPatch } from '@/services/SkillApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const patch = route.params.patch as string;

document.title = `Patch #${patch} skills - ROVis`;

const {
	state,
	total,
	getPatchItems,
} = SkillApi.use();

const list = ref<SkillPatch[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getPatchItems(patch, page);
}

loadPage(1);

const newSkills = computed(() => list.value.filter((v) => v.From === null));
const updatedSkills = computed(() => list.value.filter((v) => v.From !== null && v.To !== null));
const deletedSkills = computed(() => list.value.filter((v) => v.From !== null && v.To === null));
</script>

<template>
	<ListingBase
		:title="`Skills in patch ${patch}`"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<h4>New skills</h4>
		<BsAccordion v-if="newSkills.length > 0">
			<BsAccordionItem
				v-for="(val) in newSkills"
				:key="val.To?.Data.SkillID"
				:title="`#${val.To?.Data.SkillID} - ${val?.To?.Data.Name}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.SkillHistory, params: { skillId: val.To?.Data.SkillID ?? val.From?.Data.SkillID } }"
						target="_blank"
					>
						View Skill history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<SkillCompareTable :current="val.To?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No new skills in this page</p>

		<h4 class="mt-3">Updated skills</h4>
		<BsAccordion v-if="updatedSkills.length > 0">
			<BsAccordionItem
				v-for="(val) in updatedSkills"
				:key="val.To?.Data.SkillID ?? val.From?.Data.SkillID"
				:title="`#${val.To?.Data.SkillID ?? val.From?.Data.SkillID} - ${val?.To?.Data.Name ?? val?.From?.Data.Name}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.SkillHistory, params: { skillId: val.To?.Data.SkillID ?? val.From?.Data.SkillID } }"
						target="_blank"
					>
						View Skill history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<SkillCompareTable :current="val.To?.Data" :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No updated skills in this page</p>

		<h4 class="mt-3">Deleted skills</h4>
		<BsAccordion v-if="deletedSkills.length > 0">
			<BsAccordionItem
				v-for="(val) in deletedSkills"
				:key="val.From?.Data.SkillID"
				:title="`#${val.From?.Data.SkillID} - ${val?.From?.Data.Name}`"
			>
				<p>
					<strong>Last update:</strong> {{ val.LastUpdate }}
					<BsLink
						:to="{ name: RouteName.SkillHistory, params: { skillId: val.To?.Data.SkillID ?? val.From?.Data.SkillID } }"
						target="_blank"
					>
						View Skill history <BIconBoxArrowUpRight />
					</BsLink>
				</p>

				<SkillCompareTable :previous="val.From?.Data" />
			</BsAccordionItem>
		</BsAccordion>
		<p v-else>No deleted skills in this page</p>
	</ListingBase>
</template>
