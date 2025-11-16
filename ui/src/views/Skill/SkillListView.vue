<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import type { MinSkill } from '@/models/Skill';
import { RouteName } from '@/router/RouteName';
import { SkillApi } from '@/services/SkillApi';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { ref } from 'vue';

document.title = `Skill List - ROVis`;

const {
	state,
	total,
	getItems,
} = SkillApi.use();

const list = ref<MinSkill[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getItems(page);
}

loadPage(1);
</script>

<template>
	<ListingBase
		title="Skills"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsListGroup :flush="true">
			<BsListGroupItem
				v-for="(skill) in list"
				:key="`${skill.SkillID}-${skill.LastUpdate}`"
			>
				<BsLink
					:to="{ name: RouteName.SkillHistory, params: { skillId: skill.SkillID } }"
					target="_blank"
				>
					#{{ skill.SkillID }} - {{ skill.Constant }} [{{ skill.Name }}] (Last updated: {{ skill.LastUpdate }})
					<BIconBoxArrowUpRight />
				</BsLink>
			</BsListGroupItem>
		</BsListGroup>
	</ListingBase>
</template>
