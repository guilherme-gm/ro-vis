<script setup lang="ts">
import type { ButtonVariant } from '../Button/ButtonVariant';
import BsButton from '../Button/BsButton.vue';

defineProps<{
	label: string;
	variant: ButtonVariant;
	items: { label: string; value: string }[];
}>();

const emit = defineEmits(['item-selected']);
</script>

<template>
	<div class="dropdown">
		<BsButton :label="label" :variant="variant" class="dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false" />
		<ul class="dropdown-menu dropdown-menu-end">
			<li v-for="item in items" :key="item.value">
				<button class="dropdown-item" @click="emit('item-selected', item.value)">{{ item.label }}</button>
			</li>
		</ul>
	</div>
</template>
