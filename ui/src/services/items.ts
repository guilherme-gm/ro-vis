import { ref } from "vue";
import { useApi } from "./api";
import { LoadState } from "./LoadState";
import type { Item, MinItem } from "@/models/Item";

export type PatchItem = {
	previous: Item | null;
	current: Item | null;
};

type GetItemsResponse = {
	total: number;
	list: MinItem[];
};

type GetItemsHistoryResponse = {
	total: number;
	list: PatchItem[];
};

const api = useApi();
const total = ref(-1);
const itemHistoryTotal = ref(-1);

const state = ref<LoadState>(LoadState.None);

async function getItems(page: number): Promise<GetItemsResponse['list']> {
	try {
		state.value = LoadState.Loading;

		const items = await api.get<GetItemsResponse>('items/', { start: ((page - 1) * 100) });

		total.value = items.total;

		state.value = LoadState.Loaded;

		return items.list;
	} catch (error) {
		state.value = LoadState.Error;
	}

	return [];
}


async function getItemHistory(itemId: string, page: number): Promise<PatchItem[]> {
	try {
		state.value = LoadState.Loading;

		const items = await api.get<GetItemsHistoryResponse>('items/item.php', { itemId, start: ((page - 1) * 100) });

		itemHistoryTotal.value = items.total;

		state.value = LoadState.Loaded;

		return items.list;
	} catch (error) {
		state.value = LoadState.Error;
	}

	return [];
}

async function getPatchItems(patch: string, page: number): Promise<PatchItem[]> {
	try {
		state.value = LoadState.Loading;

		const items = await api.get<GetItemsHistoryResponse>('items/patch.php', { patch, start: ((page - 1) * 100) });

		total.value = items.total;

		state.value = LoadState.Loaded;

		return items.list;
	} catch (error) {
		state.value = LoadState.Error;
	}

	return [];
}

export function useItems() {
	return {
		state,
		total,
		itemHistoryTotal,
		getItems,
		getItemHistory,
		getPatchItems,
	};
}
