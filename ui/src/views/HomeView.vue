<script setup lang="ts">
import BsLink from '@/components/bootstrap/BsLink.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import { RouteName } from '@/router/RouteName';
import { useServerStore } from '@/stores/server';

document.title = "ROVis";

const serverStore = useServerStore();

</script>

<template>
	<BsListGroup>
		<BsListGroupItem class="d-flex justify-content-between align-items-start">
			<div class="ms-2 me-auto">
				<div class="fw-bold"><BsLink :to="{ name: RouteName.Updates }">Updates</BsLink></div>
				See details of each update, along with the related patches and links for details
			</div>
		</BsListGroupItem>
		<BsListGroupItem class="d-flex justify-content-between align-items-start">
			<div class="ms-2 me-auto">
				<div class="fw-bold"><BsLink :to="{ name: RouteName.ItemList }">Items</BsLink></div>
				Item information changes, like names, descriptions, slots, etc.
			</div>
		</BsListGroupItem>
		<BsListGroupItem class="d-flex justify-content-between align-items-start">
			<div class="ms-2 me-auto">
				<div class="fw-bold"><BsLink :to="{ name: RouteName.QuestList }">Quest Log</BsLink></div>
				Quest log updates, for actual quest entries ("Recommended Quests" not included)
			</div>
		</BsListGroupItem>
		<BsListGroupItem
			v-if="serverStore.currentServer === 'latam'"
			class="d-flex justify-content-between align-items-start"
		>
			<div class="ms-2 me-auto">
				<div class="fw-bold"><BsLink :to="{ name: RouteName.I18nList }">i18n</BsLink></div>
				Internationalization text entries and translations
			</div>
		</BsListGroupItem>
		<!-- <BsListGroupItem class="d-flex justify-content-between align-items-start">
			<div class="ms-2 me-auto">
				<div class="fw-bold"><BsLink :to="{ name: RouteName.StateList }">States</BsLink></div>
				States, like buffs and debuffs
			</div>
		</BsListGroupItem> -->
	</BsListGroup>
</template>
