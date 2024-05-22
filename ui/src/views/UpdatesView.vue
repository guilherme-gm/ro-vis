<script setup lang="ts">
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import { useUpdates } from '@/services/updates';

const { getUpdates, isLoading, updateList } = useUpdates();
getUpdates();
</script>

<template>
	<h2>Updates History</h2>
	<div v-if="!isLoading">
		<BsAccordion>
			<BsAccordionItem
				v-for="(val) in updateList"
				:key="val.id"
				:title="val.id"
			>
				<h5>Changed files:</h5>
				<BsListGroup :flush="true">
					<BsListGroupItem
						v-for="(update) in val.updates"
						:key="`${val.id}-${update.file}`"
					>
						{{ update.file }}
					</BsListGroupItem>
				</BsListGroup>
			</BsAccordionItem>
		</BsAccordion>
	</div>
</template>
