<script setup lang="ts">
import { nextTick, ref } from 'vue';

const isCollapsing = ref(false);

async function beforeLeave(el: Element): Promise<void> {
	const htmlEl = el as HTMLElement;
	isCollapsing.value = true;
	htmlEl.style['height'] = `${el.getBoundingClientRect().height}px`;
	await nextTick();
	htmlEl.style['height'] = '';
}

async function beforeEnter(el: Element): Promise<void> {
	const htmlEl = el as HTMLElement;
	isCollapsing.value = true;
	htmlEl.style['height'] = '';
	await nextTick();
	htmlEl.style['height'] = `${htmlEl.scrollHeight}px`;
}

function afterEnter(el: Element): void {
	const htmlEl = el as HTMLElement;
	isCollapsing.value = false;
	htmlEl.style['height'] = '';
}
</script>

<template>
	<Transition
		leave-active-class="collapsing"
		enter-active-class="collapsing"
		@before-leave="beforeLeave"
		@after-leave="isCollapsing = false"
		@before-enter="beforeEnter"
		@after-enter="afterEnter"
	>
		<slot></slot>
	</Transition>
</template>
