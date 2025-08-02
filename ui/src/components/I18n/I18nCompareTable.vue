<script setup lang="ts">
import DiffedValue from "@/components/DiffedValue.vue";
import { computed } from "vue";
import type { I18n } from "@/models/I18n";

const props = defineProps<{
	previous?: I18n | null;
	current?: I18n | null;
}>();

const fields = [
	{ key: 'I18nId', label: 'ID' },
	{ key: 'ContainerFile', label: 'Container File' },
	{ key: 'EnText', label: 'English Text' },
	{ key: 'PtBrText', label: 'Portuguese Text' },
	{ key: 'FileVersion', label: 'File Version' },
	{ key: 'Active', label: 'Active' },
	{ key: 'Deleted', label: 'Deleted' },
];

const showNew = computed(() => props.current);
const showPrevious = computed(() => props.previous);
const showDiff = computed(() => props.previous && props.current);
</script>

<template>
	<table class="table table-striped table-sm">
		<tbody>
			<tr>
				<th>Field</th>
				<th v-if="showNew">New</th>
				<th v-if="showPrevious">Previous</th>
				<th v-if="showDiff">Diff</th>
			</tr>
			<tr v-for="field in fields" :key="field.key">
				<th>{{ field.label }}</th>
				<td v-if="showNew">
					<pre class="pre-preserve">{{ current?.[field.key as keyof I18n] ?? "-" }}</pre>
				</td>
				<td v-if="showPrevious">
					<pre class="pre-preserve">{{ previous?.[field.key as keyof I18n] ?? "-" }}</pre>
				</td>
				<td v-if="showDiff">
					<DiffedValue :from="previous?.[field.key as keyof I18n]" :to="current?.[field.key as keyof I18n]" />
				</td>
			</tr>
		</tbody>
	</table>
</template>

<style scoped>
	.pre-preserve {
		white-space: pre-wrap;
		word-wrap: break-word;
		margin: 0;
		font-family: inherit;
		font-size: 1em;
		line-height: 1.5;
	}
</style>
