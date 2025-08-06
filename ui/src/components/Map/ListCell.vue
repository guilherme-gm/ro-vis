<script setup lang="ts" generic="T">
import { computed } from "vue";
import { ListFormatter } from "./formatters";
import type { Formatter } from "./formatters";

const props = defineProps<{
	value?: T[];
	formatter: Formatter<T>;
}>();

const listFormatter = ListFormatter.use(props.formatter);

const list = computed(() => listFormatter.formatList(props.value));
</script>

<template>
	<ul>
		<li v-for="(item, index) of list" :key="index">
			<slot :item="item">{{ item.stringValue }}</slot>
		</li>
	</ul>
</template>
