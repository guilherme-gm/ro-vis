<script setup lang="ts">
import { computed } from "vue";
import type { MapNpc } from "@/models/Map";
import { MapNpcFormatter } from "./formatters";

const props = defineProps<{
	value?: MapNpc[];
}>();

const formatter = MapNpcFormatter.use();

const npcList = computed(() => formatter.formatList(props.value));
</script>

<template>
	<ul>
		<li v-for="(npc, index) of npcList" :key="index">{{ npc }}</li>
	</ul>
</template>
