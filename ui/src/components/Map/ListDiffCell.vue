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
			v-for="(item, index) of diffFormatter.formatList()"
			:key="index"
		>
			<slot :item="item" :kind="item.diffType">
				<span :class="{ [`diff-${item.diffType}`]: true }">{{ item.stringValue }}</span>
			</slot>
		</li>
	</ul>
</template>
