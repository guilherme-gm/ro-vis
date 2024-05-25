<script setup lang="ts">
import DiffedValue from "@/components/DiffedValue.vue";
import type { Item } from "@/models/Item";
import { computed, ref } from "vue";
import MoveInfoCell from "./MoveInfoCell.vue";
import MoveInfoDiffCell from "./MoveInfoDiffCell.vue";

const props = defineProps<{
	previous?: Item | null;
	current?: Item | null;
}>();

const fields = ref<[string, keyof Item][]>([
	['Identified name', 'IdentifiedName'],
	['Identified description', 'IdentifiedDescription'],
	['Identified sprite', 'IdentifiedSprite'],
	['Unidentified name', 'UnidentifiedName'],
	['Unidentified description', 'UnidentifiedDescription'],
	['Unidentified sprite', 'UnidentifiedSprite'],
	['Slot count', 'SlotCount'],
	['Is book', 'IsBook'],
	['Can use buying store', 'CanUseBuyingStore'],
	['Card prefix', 'CardPrefix'],
	['Card postfix', 'CardPostfix'],
	['Card illustration', 'CardIllustration'],
	['Class num', 'ClassNum'],
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
				<th>Move Info</th>
				<td v-if="showNew" ><MoveInfoCell :info="current?.MoveInfo ?? null" /></td>
				<td v-if="showPrevious"><MoveInfoCell :info="current?.MoveInfo ?? null" /></td>
				<td v-if="showDiff" ><MoveInfoDiffCell :from="previous?.MoveInfo ?? null" :to="current?.MoveInfo ?? null" /></td>
			</tr>
		</tbody>
	</table>
</template>
