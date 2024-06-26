<script setup lang="ts">
import type { State } from '@/models/State';
import ColorBox from '@/components/ColorBox.vue';
import ATooltip from '../ATooltip.vue';
import { BIconClock } from 'bootstrap-icons-vue';
import { computed } from 'vue';
import DiffedValue from '../DiffedValue.vue';

const props = defineProps<{
	from: State | null | undefined;
	to: State | null | undefined;
}>();

type DiffedLine = {
	type: 'simple-text';
	text: string;
	color?: number[];
	hasTimer?: boolean;
	class: 'added' | 'removed' | 'no-change';
} | {
	type: 'diff'
	didColorChange: boolean;
	fromText: string;
	fromColor: number[] | null;
	fromTimer: boolean;
	toText: string;
	toColor: number[] | null;
	toTimer: boolean;
};

function areArraysEqual(arr1: number[] | null, arr2: number[] | null) {
	if (!arr1 && !arr2) {
		return true;
	}

	if (!arr1 || !arr2) {
		return false;
	}

	if (arr1.length !== arr2.length) {
		return false;
	}

	for (let i = 0; i < arr1.length; i++) {
		if (arr1[i] !== arr2[i]) {
			return false;
		}
	}

	return true;
}

const diffedLines = computed((): DiffedLine[] => {
	if ((props.from?.Description.length ?? 0) === 0 && (props.to?.Description.length ?? 0) !== 0) {
		return [{ type: 'simple-text', text: 'New description', class: 'added' }];
	}

	if ((props.from?.Description.length ?? 0) !== 0 && (props.to?.Description.length ?? 0) === 0) {
		return [{ type: 'simple-text', text: 'Description removed', class: 'removed' }];
	}

	if ((props.from?.Description.length ?? 0) === 0 && (props.to?.Description.length ?? 0) === 0) {
		return [{ type: 'simple-text', text: 'N/A', class: 'no-change' }];
	}

	const diff: DiffedLine[] = [];

	const commonLen = Math.min(props.from?.Description.length ?? 0, props.to?.Description.length ?? 0);
	let i;
	for (i = 0; i < commonLen; i++) {
		const fromLine = props.from?.Description[i];
		const toLine = props.to?.Description[i];

		diff.push({
			type: 'diff',
			didColorChange: !areArraysEqual(fromLine?.Color ?? null, toLine?.Color ?? null),
			fromText: fromLine?.Text ?? '',
			fromColor: fromLine?.Color ?? null,
			fromTimer: props.from?.TimeStrLineNum === i || false,
			toText: toLine?.Text ?? '',
			toColor: toLine?.Color ?? null,
			toTimer: props.to?.TimeStrLineNum === i || false,
		});
	}

	for (let j = i; j < (props.from?.Description.length ?? 0); j++) {
		diff.push({ type: 'simple-text', text: props.from?.Description[j]?.Text ?? '', class: 'removed' });
	}

	for (let j = i; j < (props.to?.Description.length ?? 0); j++) {
		diff.push({ type: 'simple-text', text: props.to?.Description[j]?.Text ?? '', class: 'added' });
	}

	return diff;
});
</script>

<template>
	<div v-for="(line, idx) in diffedLines" :key="idx" class="d-flex flex-column gap-3">
		<span v-if="line.type === 'simple-text'" :class="line.class" class="d-flex flex-row gap-3">
			<ColorBox v-if="line.color" :color="line.color" :default-color="[109, 109, 109]" />
			{{ line.text }}
			<ATooltip v-if="line.hasTimer" text="This line shows the remaining time">
				<BIconClock />
			</ATooltip>
		</span>
		<span v-else class="d-flex flex-row gap-3">
			<span v-if="line.didColorChange">
				<ColorBox :color="line.fromColor" :default-color="[109, 109, 109]" />
				--&gt;
				<ColorBox v-if="line.toColor" :color="line.toColor" :default-color="[109, 109, 109]" />
			</span>
			<span v-else-if="line.fromColor">
				<ColorBox :color="line.fromColor" :default-color="[109, 109, 109]" />
			</span>

			<DiffedValue :from="line.fromText" :to="line.toText" />

			<span v-if="line.fromTimer !== line.toTimer">
				<ATooltip v-if="line.fromTimer" text="This line no longer shows remaining time" class="removed">
					<BIconClock />
				</ATooltip>
				<ATooltip v-if="line.toTimer" text="This line now shows the remaining time" class="added">
					<BIconClock />
				</ATooltip>
			</span>
			<span v-else-if="line.fromTimer">
				<ATooltip v-if="line.toTimer" text="This line continues showing the remaining time" class="no-change">
					<BIconClock />
				</ATooltip>
			</span>
		</span>
	</div>
</template>

<style scoped>
	.removed {
		color: rgb(236, 85, 85);
	}

	.added {
		color: rgb(89, 189, 114);
	}

	.no-change {
		color: grey;
	}
</style>
