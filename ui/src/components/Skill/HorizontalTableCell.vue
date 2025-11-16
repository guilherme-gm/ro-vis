<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	value: unknown[];
	rows: {
		label: string;
		value: (it: unknown, idx: number) => string;
	}[];
}>();

const formattedValues = computed(() => {
	return props.rows.map((rowConfig) => {
		return props.value.map((val, idx) => {
			return rowConfig.value(val, idx);
		});
	});
});
</script>

<template>
	<table v-if="value?.length > 0" class="table table-bordered w-auto">
		<tr v-for="(it, idx) of rows" :key="idx">
			<th class="p-2">{{ it.label }}</th>
			<td v-for="column in formattedValues[idx]" :key="column" class="p-2">{{ column }}</td>
		</tr>
	</table>
	<span v-else class="text-muted">-</span>
</template>
