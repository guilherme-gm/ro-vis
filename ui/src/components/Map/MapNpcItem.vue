<script setup lang="ts">
import type { MapNpc } from "@/models/Map";
import LocalizableString from "../LocalizableString/LocalizableString.vue";
import { computed } from "vue";

const props = defineProps({
	npc: {
		type: Object as () => MapNpc,
		required: true,
	},
	kind: {
		type: String,
		validator: (value: string) => ['added', 'removed', 'unchanged', 'none'].includes(value),
		default: 'none',
	},
});

const hiddenPart = computed(() => props.npc.Name2 ? ` (#${props.npc.Name2})` : '');
</script>

<template>
	<span :class="{ [`diff-${kind}`]: true }">[{{ npc.Location.X }}, {{ npc.Location.Y }}] <LocalizableString :string="npc.Name1" /> {{ hiddenPart }}</span>
</template>
