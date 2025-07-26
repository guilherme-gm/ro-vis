<script setup lang="ts">
import { computed } from 'vue';
import { ButtonVariant } from './ButtonVariant';

const props = defineProps<{
	label: string;
	variant?: ButtonVariant;
}>();

const emit = defineEmits(['click']);

const buttonClass = computed(() => ({
	btn: true,
	[(props.variant || ButtonVariant.Primary)]: true,
}));
</script>

<template>
	<button :class="buttonClass" @click="emit('click')">
		<slot name="icon"></slot>
		{{ label }}</button>
</template>
