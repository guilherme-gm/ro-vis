<script setup lang="ts">
import type { State } from '@/models/State';
import ColorBox from '@/components/ColorBox.vue';
import ATooltip from '../ATooltip.vue';
import { BIconClock } from 'bootstrap-icons-vue';

defineProps<{
	state: State | null | undefined;
}>();
</script>

<template>
	<span v-if="!state || state.Description.length === 0">N/A</span>
	<div v-else class="d-flex flex-column gap-3">
		<span v-for="(line, idx) in state.Description" :key="idx">
			<ColorBox :color="line.Color" :default-color="[109, 109, 109]" />
			{{ line.Text }}
			<ATooltip v-if="idx === state.TimeStrLineNum" text="This line shows the remaining time">
				<BIconClock />
			</ATooltip>
		</span>
	</div>
</template>
