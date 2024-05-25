import { ref } from "vue";
import { useApi } from "./api";
import { LoadState } from "./LoadState";
import type { Item } from "@/models/Item";

export type PatchItem = {
	previous: Item | null;
	current: Item | null;
};

type GetItemsResponse = {
	total: number;
	list: PatchItem[];
};

const api = useApi();
const total = ref(-1);

const state = ref<LoadState>(LoadState.None);

/*
async function getItems(page: number): Promise<Item[]> {
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
*/

async function getPatchItems(patch: string, page: number): Promise<PatchItem[]> {
	try {
		state.value = LoadState.Loading;

		const items = await api.get<GetItemsResponse>('items/patch.php', { patch, start: ((page - 1) * 100) });

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
		// getItems,
		getPatchItems,
	};
}
