<script setup lang="ts">
import DiffedValue from "@/components/DiffedValue.vue";
import { computed, ref } from "vue";
import type { State } from "@/models/State";
import DescriptionCell from "./DescriptionCell.vue";
import DescriptionDiffCell from "./DescriptionDiffCell.vue";

const props = defineProps<{
	previous?: State | null;
	current?: State | null;
}>();

const fields = ref<[string, keyof State][]>([
	// ['Description', 'Description'], // custom logic
	['Has time limit', 'HasTimeLimit'],
	['Time line index', 'TimeStrLineNum'],
	['Has effect image', 'HasEffectImage'],
	['Icon image', 'IconImage'],
	['Icon priority', 'IconPriority'],
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
			<tr>
				<th>Description</th>
				<td v-if="showNew"><DescriptionCell :state="current" /></td>
				<td v-if="showPrevious"><DescriptionCell :state="previous" /></td>
				<td v-if="showDiff"><DescriptionDiffCell :from="previous" :to="current" /></td>
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
