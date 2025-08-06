<script setup lang="ts">
import { computed } from "vue";
import type { MapNpc } from "@/models/Map";
import { ListFormatter, MapNpcFormatter } from "./formatters";

const props = defineProps<{
	value?: MapNpc[];
}>();

const formatter = MapNpcFormatter.use();
const listFormatter = ListFormatter.use(formatter);

const npcList = computed(() => listFormatter.formatList(props.value));
</script>

<template>
	<ul>
		<li v-for="(npc, index) of npcList" :key="index">{{ npc }}</li>
	</ul>
</template>
