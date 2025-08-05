<script setup lang="ts">
import DiffedValue from "@/components/DiffedValue.vue";
import { computed, ref } from "vue";
import type { Map } from "@/models/Map";

const props = defineProps<{
	previous?: Map | null;
	current?: Map | null;
}>();

const fields = ref<[string, keyof Map][]>([
	['Name', 'Name'],
	['SpecialCode', 'SpecialCode'],
	['Mp3Name', 'Mp3Name'],
	['Npcs', 'Npcs'],
	['Warps', 'Warps'],
	['Spawns', 'Spawns'],
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
