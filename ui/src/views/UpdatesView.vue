<script setup lang="ts">
import ListingBase from '@/components/ListingBase.vue';
import BsAccordion from '@/components/bootstrap/Accordion/BsAccordion.vue';
import BsAccordionItem from '@/components/bootstrap/Accordion/BsAccordionItem.vue';
import BsLink from '@/components/bootstrap/BsLink.vue';
import BsListGroup from '@/components/bootstrap/ListGroup/BsListGroup.vue';
import BsListGroupItem from '@/components/bootstrap/ListGroup/BsListGroupItem.vue';
import type { Update } from '@/models/Update';
import { RouteName } from '@/router/RouteName';
import { useUpdates } from '@/services/updates';
import { BIconBoxArrowUpRight } from 'bootstrap-icons-vue';
import { ref } from 'vue';

document.title = "Update history - ROVis";

const {
	state,
	total,
	getUpdates,
} = useUpdates();

const list = ref<Update[]>([]);
const currentPage = ref(0);

async function loadPage(page: number): Promise<void> {
	currentPage.value = page;
	list.value = await getUpdates(page);
}

loadPage(1);

type TrackInfo = {
	route: RouteName;
	label: string;
};

const itemTrackInfo: TrackInfo = {
	route: RouteName.ItemPatch,
	label: 'Item',
};

const questTrackInfo: TrackInfo = {
	route: RouteName.QuestPatch,
	label: 'Quest',
};

const i18nTrackInfo: TrackInfo = {
	route: RouteName.I18nPatch,
	label: 'I18n',
};

const fileTrackInfo = new Map<string, { route: RouteName; label: string; }>([
	// ======= Item
	["data/bookitemnametable.txt", itemTrackInfo],
	["data/buyingstoreitemlist.txt", itemTrackInfo],
	["data/cardpostfixnametable.txt", itemTrackInfo],
	["data/cardprefixnametable.txt", itemTrackInfo],
	["data/idnum2itemdesctable.txt", itemTrackInfo],
	["data/idnum2itemdisplaynametable.txt", itemTrackInfo],
	["data/idnum2itemresnametable.txt", itemTrackInfo],
	["data/itemslotcounttable.txt", itemTrackInfo],
	["data/num2cardillustnametable.txt", itemTrackInfo],
	["data/num2itemdesctable.txt", itemTrackInfo],
	["data/num2itemdisplaynametable.txt", itemTrackInfo],
	["data/num2itemresnametable.txt", itemTrackInfo],
	// new on v2
	["System/iteminfo.lub", itemTrackInfo],
	// new on v3
	["data/itemmoveinfov5.txt", itemTrackInfo],

	// ====== Quest
	["data/questid2display.txt", questTrackInfo],
	["system/ongoingquestinfolist_true.lub", questTrackInfo],

	// ====== I18n
	// ["data/i18n/sc/(.*?).csv", i18nTrackInfo],
	["data/i18n/sc/sc.json", i18nTrackInfo]
]);

function isTrackedFile(fileName: string): boolean {
	if (/^data\/i18n\/sc\/(.*?)\.csv$/i.test(fileName.toLocaleLowerCase())) {
		return true;
	}

	return fileTrackInfo.has(fileName.toLocaleLowerCase());
}

function patchRoute(fileName: string): RouteName {
	if (/^data\/i18n\/sc\/(.*?)\.csv$/i.test(fileName.toLocaleLowerCase())) {
		return i18nTrackInfo.route;
	}

	return fileTrackInfo.get(fileName.toLocaleLowerCase())!.route;
}

function trackLabel(fileName: string): string {
	if (/^data\/i18n\/sc\/(.*?)\.csv$/i.test(fileName.toLocaleLowerCase())) {
		return i18nTrackInfo.label;
	}

	return fileTrackInfo.get(fileName.toLocaleLowerCase())!.label;
}
</script>

<template>
	<ListingBase
		title="Update history"
		:total="total"
		:state="state"
		:current-page="currentPage"
		@page-changed="loadPage"
	>
		<BsAccordion>
			<BsAccordionItem
				v-for="(val) in list"
				:key="val.date"
				:title="val.date"
			>
				<h5>Changed files:</h5>
				<BsListGroup :flush="true">
					<BsListGroupItem
						v-for="(update) in val.changes"
						:key="`${val.date}-${update.file}`"
					>
						<BsLink
							v-if="isTrackedFile(update.file)"
							:to="{ name: patchRoute(update.file), params: { patch: val.date } }"
							target="_blank"
						>
							({{ trackLabel(update.file) }}) {{ update.file }}&nbsp;&nbsp;
							<BIconBoxArrowUpRight />
						</BsLink>
						<span v-else>{{ update.file }}</span>
					</BsListGroupItem>
				</BsListGroup>
			</BsAccordionItem>
		</BsAccordion>
	</ListingBase>
</template>
