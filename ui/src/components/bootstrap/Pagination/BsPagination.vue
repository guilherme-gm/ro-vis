<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	numPages: number;
	currentPage: number;
	description: string;
	center?: boolean;
	groupSize?: number;
}>();

const emit = defineEmits(['changed']);

const pageList = computed((): { label: string; value: number; disabled: boolean; }[] => {
	const pages: { label: string; value: number; disabled: boolean }[] = [];
	if (!props.groupSize || props.numPages < props.groupSize) {
		for (let i = 1; i <= props.numPages; i++) {
			pages.push({ label: i.toString(), value: i, disabled: false });
		}

		return pages;
	}

	const currentIndex = props.currentPage - 1;
	const currentSlice = Math.floor(currentIndex / props.groupSize);
	const sliceStartIndex = currentSlice * props.groupSize;
	const sliceEndIndex = ((currentSlice + 1) * props.groupSize) - 1;

	console.log(currentIndex, currentSlice, sliceStartIndex, sliceEndIndex);

	pages.push({
		label: 'Previous',
		value: sliceStartIndex,
		disabled: sliceStartIndex === 0,
	});

	for (let pageNum = sliceStartIndex + 1; ((pageNum <= sliceEndIndex + 1) && pageNum <= props.numPages); pageNum += 1) {
		pages.push({
			label: pageNum.toString(),
			value: pageNum,
			disabled: false,
		});
	}

	pages.push({
		label: 'Next',
		value: sliceEndIndex + 2,
		disabled: props.numPages < (sliceEndIndex + 1),
	});

	return pages;
});
</script>

<template>
	<nav :aria-label="description">
		<ul class="pagination" :class="{ 'justify-content-center': center }">
			<li
				v-for="(page) in pageList"
				:key="`page-${page.value}`"
				class="page-item"
				:class="{
					active: (page.value === currentPage),
					disabled: page.disabled,
				}"
			>
				<button class="page-link" @click="emit('changed', page.value)">{{ page.label }}</button>
			</li>
		</ul>
	</nav>
</template>
