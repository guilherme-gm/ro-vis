<script setup lang="ts" generic="T">
import { computed } from "vue";
import { ListDiffFormatter, ListFormatter } from "./formatters";
import { ListDiffer } from "./differs";
import type { Formatter } from "./formatters";

const props = defineProps<{
	formatter: Formatter<T>;
	from?: T[];
	to?: T[];
}>();

const differ = ListDiffer.use(computed(() => props.from ?? []), computed(() => props.to ?? []), props.formatter);
const listFormatter = ListFormatter.use(props.formatter);
const diffFormatter = ListDiffFormatter.use(differ, listFormatter, props.formatter);
</script>

<template>
	<ul>
		<li
			v-for="(npc, index) of diffFormatter.formatList()"
			:key="index"
			:class="{ [`diff-${npc.diffType}`]: true }"
		>
			{{ npc.value }}
		</li>
	</ul>
</template>
