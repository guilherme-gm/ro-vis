<script setup lang="ts">
import { computed } from "vue";
import type { MapNpc } from "@/models/Map";
import { MapNpcDiffFormatter } from "./formatters";
import { MapNpcDiffer } from "./differs";

const props = defineProps<{
	from?: MapNpc[];
	to?: MapNpc[];
}>();

const differ = MapNpcDiffer.use(computed(() => props.from ?? []), computed(() => props.to ?? []));
const diffFormatter = MapNpcDiffFormatter.use(differ);
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
