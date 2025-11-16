<script setup lang="ts">
import DiffedValue from "@/components/DiffedValue.vue";
import { computed, ref, type Component } from "vue";
import type { Skill } from "@/models/Skill";
import HorizontalTableCell from "./HorizontalTableCell.vue";
import HorizontalTableDiffCell from "./HorizontalTableDiffCell.vue";

const props = defineProps<{
	previous?: Skill | null;
	current?: Skill | null;
}>();

const fields = ref<{
	label: string;
	key: keyof Skill;
	component?: Component | null;
	diffComponent?: Component | null;
	options?: Record<string, unknown>;
}[]>([
	{ label: 'Constant', key: 'Constant' },
	{ label: 'Name', key: 'Name' },
	{ label: 'Description', key: 'Description' },
	{ label: 'Max Level', key: 'MaxLevel' },
	{ label: 'Is Passive', key: 'IsPassive' },
	{
		label: 'SP Cost',
		key: 'SpCost',
		component: HorizontalTableCell,
		diffComponent: HorizontalTableDiffCell,
		options: {
			rows: [
				{ label: 'Lv', value: (it: Skill['SpCost'][0], idx: number) => idx + 1 },
				{ label: 'SP', value: (it: Skill['SpCost'][0]) => it },
			]
		},
	},
	{
		label: 'AP Cost',
		key: 'ApCost',
		component: HorizontalTableCell,
		diffComponent: HorizontalTableDiffCell,
		options: {
			rows: [
				{ label: 'Lv', value: (it: Skill['ApCost'][0], idx: number) => idx + 1 },
				{ label: 'AP', value: (it: Skill['ApCost'][0]) => it },
			]
		},
	},
	{ label: 'Can Select Level', key: 'CanSelectLevel' },
	{
		label: 'Attack Range',
		key: 'AttackRange',
		component: HorizontalTableCell,
		diffComponent: HorizontalTableDiffCell,
		options: {
			rows: [
				{ label: 'Lv', value: (it: Skill['AttackRange'][0], idx: number) => idx + 1 },
				{ label: 'Range', value: (it: Skill['AttackRange'][0]) => it },
			]
		},
	},
	{ label: 'Required Skills', key: 'RequiredSkills' },
	{ label: 'Job Required Skills', key: 'JobRequiredSkills' },
	{
		label: 'Skill Scale',
		key: 'SkillScale',
		component: HorizontalTableCell,
		diffComponent: HorizontalTableDiffCell,
		options: {
			rows: [
				{ label: 'Lv', value: (it: Skill['SkillScale'][0]) => it.Level },
				{ label: 'X', value: (it: Skill['SkillScale'][0]) => it.X },
				{ label: 'Y', value: (it: Skill['SkillScale'][0]) => it.Y },
			]
		},
	},
	{ label: 'Cast Flags', key: 'CastFlags' },
	{
		label: 'Cast Fixed Delay',
		key: 'CastFixedDelay',
		component: HorizontalTableCell,
		diffComponent: HorizontalTableDiffCell,
		options: {
			rows: [
				{ label: 'Lv', value: (it: Skill['CastFixedDelay'][0], idx: number) => idx + 1 },
				{ label: 'Delay', value: (it: Skill['CastFixedDelay'][0]) => it }
			]
		}
	},
	{
		label: 'Cast Stat Delay',
		key: 'CastStatDelay',
		component: HorizontalTableCell,
		diffComponent: HorizontalTableDiffCell,
		options: {
			rows: [
				{ label: 'Lv', value: (it: Skill['CastStatDelay'][0], idx: number) => idx + 1 },
				{ label: 'Delay', value: (it: Skill['CastStatDelay'][0]) => it }
			]
		}
	},
	{
		label: 'Single Post Delay',
		key: 'SinglePostDelay',
		component: HorizontalTableCell,
		diffComponent: HorizontalTableDiffCell,
		options: {
			rows: [
				{ label: 'Lv', value: (it: Skill['SinglePostDelay'][0], idx: number) => idx + 1 },
				{ label: 'Delay', value: (it: Skill['SinglePostDelay'][0]) => it }
			]
		}
	},
	{
		label: 'Global Post Delay',
		key: 'GlobalPostDelay',
		component: HorizontalTableCell,
		diffComponent: HorizontalTableDiffCell,
		options: {
			rows: [
				{ label: 'Lv', value: (it: Skill['GlobalPostDelay'][0], idx: number) => idx + 1 },
				{ label: 'Delay', value: (it: Skill['GlobalPostDelay'][0]) => it }
			]
		}
	},
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
			<tr v-for="(info) of fields" :key="info.key">
				<th>{{ info.label }}</th>
				<td v-if="showNew">
					<component v-if="info.component" :is="info.component" v-bind="info.options"
						:value="current?.[info.key]" />
					<pre v-else class="pre-preserve">{{ current?.[info.key] ?? "-" }}</pre>
				</td>
				<td v-if="showPrevious">
					<component v-if="info.component" :is="info.component" v-bind="info.options"
						:value="previous?.[info.key]" />
					<pre v-else class="pre-preserve">{{ previous?.[info.key] ?? "-" }}</pre>
				</td>
				<td v-if="showDiff">
					<component v-if="info.diffComponent" :is="info.diffComponent" v-bind="info.options"
						:from="previous![info.key]" :to="current![info.key]" />
					<DiffedValue v-else :from="previous![info.key]" :to="current![info.key]" />
				</td>
			</tr>
		</tbody>
	</table>
</template>
