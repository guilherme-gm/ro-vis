<script setup lang="ts">
import DiffedValue from "@/components/DiffedValue.vue";
import { computed, ref } from "vue";
import type { Quest } from "@/models/Quest";

const props = defineProps<{
	previous?: Quest | null;
	current?: Quest | null;
}>();

const fields = ref<[string, keyof Quest][]>([
	['Title', 'Title'],
	['Description', 'Description'],
	['Summary', 'Summary'],
	['OldImage', 'OldImage'],
	['IconName', 'IconName'],
	['NpcSpr', 'NpcSpr'],
	['NpcNavi', 'NpcNavi'],
	['NpcPosX', 'NpcPosX'],
	['NpcPosY', 'NpcPosY'],
	['CoolTimeQuest', 'CoolTimeQuest'],
	['RewardEXP', 'RewardEXP'],
	['RewardJEXP', 'RewardJEXP'],
	['RewardItemList', 'RewardItemList'],
]);

const showNew = computed(() => props.current);
const showPrevious = computed(() => props.previous);
const showDiff = computed(() => props.previous && props.current);
</script>


<template>
	<table class="table table-striped table-sm">
		<tbody>
			<tr>
				<th>Info</th>
				<th v-if="showNew">New</th>
				<th v-if="showPrevious">Previous</th>
				<th v-if="showDiff">Diff</th>
			</tr>
			<tr v-for="(info) of fields" :key="info[1]">
				<th>{{ info[0] }}</th>
				<td v-if="showNew"><pre class="pre-preserve">{{ current?.[info[1]] ?? "-" }}</pre></td>
				<td v-if="showPrevious"><pre class="pre-preserve">{{ previous?.[info[1]] ?? "-" }}</pre></td>
				<td v-if="showDiff"><DiffedValue :from="previous![info[1]]" :to="current![info[1]]" /></td>
			</tr>
		</tbody>
	</table>
</template>
