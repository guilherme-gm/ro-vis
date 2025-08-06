<script setup lang="ts">
import type { MapWarp } from "@/models/Map";
import LocalizableString from "../LocalizableString/LocalizableString.vue";
import { computed } from "vue";

const props = defineProps({
	warp: {
		type: Object as () => MapWarp,
		required: true,
	},
	kind: {
		type: String,
		validator: (value: string) => ['added', 'removed', 'unchanged', 'none'].includes(value),
		default: 'none',
	},
});

const hiddenPart = computed(() => props.warp.Name2 ? ` (#${props.warp.Name2})` : '');
</script>

<template>
	<div class="d-flex flex-column" :class="{ [`diff-${kind}`]: true }">
		<span>
			({{ warp.WarpType }}) {{ warp.From.MapId }} ({{ warp.From.X }}, {{ warp.From.Y }}) -> {{ warp.To.MapId }} ({{ warp.To.X }}, {{ warp.To.Y }})
		</span>
		<span><LocalizableString :string="warp.Name1" /> {{ hiddenPart }} ({{ warp.SpriteId }})</span>
	</div>
</template>
