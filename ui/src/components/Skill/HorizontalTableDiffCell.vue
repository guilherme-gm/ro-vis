<script setup lang="ts">
import { computed } from 'vue';
import DiffedValue from '../DiffedValue.vue';

const props = defineProps<{
	from: unknown[];
	to: unknown[];
	rows: {
		label: string;
		value: (it: unknown, idx: number) => string;
	}[];
}>();

const formattedValues = computed(() => {
	return props.rows.map((rowConfig) => {
		const slots = Math.max(props.from?.length ?? 0, props.to?.length ?? 0);
		let from = Array(slots).fill(null);
		let to = Array(slots).fill(null);

		for (let i = 0; i < slots; i++) {
			let fromVal = props.from?.[i] ?? '-';
			if (fromVal !== '-') {
				fromVal = rowConfig.value(fromVal, i);
			}
			from[i] = fromVal;

			let toVal = props.to?.[i] ?? '-';
			if (toVal !== '-') {
				toVal = rowConfig.value(toVal, i);
			}
			to[i] = toVal;
		}
		return { from, to };
	});
});

const hasValues = computed(() => formattedValues.value[0].from.length > 0);
</script>

<template>
	<table v-if="hasValues" class="table table-bordered w-auto">
		<tr v-for="(it, idx) of rows" :key="idx">
			<th class="p-2">{{ it.label }}</th>
			<td v-for="(val, colIdx) of formattedValues[idx].from" :key="colIdx" class="p-2">
				<DiffedValue :from="val" :to="formattedValues[idx].to[colIdx]" />
			</td>
		</tr>
	</table>
	<span v-else class="text-muted">-</span>
</template>
