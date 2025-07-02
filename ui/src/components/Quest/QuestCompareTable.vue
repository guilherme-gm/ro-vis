<script setup lang="ts">
import DiffedValue from "@/components/DiffedValue.vue";
import RewardItemsCell from "./RewardItemsCell.vue";
import RewardItemsDiffCell from "./RewardItemsDiffCell.vue";
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
			<tr>
				<th>RewardItemList</th>
				<td v-if="showNew"><RewardItemsCell :rewards="current?.RewardItemList ?? []" /></td>
				<td v-if="showPrevious"><RewardItemsCell :rewards="previous?.RewardItemList ?? []" /></td>
				<td v-if="showDiff"><RewardItemsDiffCell :from="previous?.RewardItemList ?? []" :to="current?.RewardItemList ?? []" /></td>
			</tr>
		</tbody>
	</table>
</template>
