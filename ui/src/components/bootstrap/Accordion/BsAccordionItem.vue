<script setup lang="ts">
import { getUuid } from '@/utils/uuid';
import BsCollapse from '@/components/bootstrap/Transitions/BsCollapse.vue'

const expanded = defineModel('expanded', { default: false });

withDefaults(
	defineProps<{
		title?: string;
		/**
		 * Simplifies the visualization by removing borders and rounded corners.
		 */
		flush?: boolean;
	}>(),
	{
		title: "",
		flush: false,
	},
);

const uuid = getUuid();
</script>

<template>
	<div class="accordion-item">
		<h2 class="accordion-header">
			<button
				type="button"
				class="accordion-button"
				:class="{ collapsed: !expanded }"
				:aria-expanded="expanded"
				:aria-controls="uuid"
				@click="expanded = !expanded"
			>
				<slot name="title">{{ title }}</slot>
			</button>
		</h2>
		<BsCollapse>
			<div v-show="expanded" :id="uuid" class="accordion-collapse" :class="{ collapse: expanded, show: expanded }">
				<div class="accordion-body">
					<slot />
				</div>
			</div>
		</BsCollapse>
	</div>
</template>
