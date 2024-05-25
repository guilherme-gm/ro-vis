<script setup lang="ts">
import { diffText } from "@/utils/diffUtils";
import { nextTick, ref, watch } from "vue";

const props = defineProps<{
	from: unknown;
	to: unknown;
}>();

const diffSlot = ref<HTMLPreElement | null>(null);

function reDiff() {
	const fromText = (props.from as string)?.toString();
	const toText = (props.to as string)?.toString();

	const fragment = diffText(fromText, toText);

	nextTick(() => {
		if (diffSlot.value !== null) {
			diffSlot.value.innerHTML = '';
			diffSlot.value.appendChild(fragment);
		}
	});
}

watch(props, () => reDiff());
reDiff();
</script>

<template>
	<pre ref="diffSlot"></pre>
</template>
