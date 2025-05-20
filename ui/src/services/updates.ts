import { ref } from "vue";
import type { Update } from "@/models/Update";
import { useApi } from "./api";
import { LoadState } from "./LoadState";

type GetUpdatesResponse = {
	total: number;
	list: Update[];
};

const api = useApi();
const total = ref(-1);

const state = ref<LoadState>(LoadState.None);

async function getUpdates(page: number): Promise<Update[]> {
	try {
		state.value = LoadState.Loading;

		const list = await api.get<GetUpdatesResponse>('updates/', { start: ((page - 1) * 100) });

		const updateList = list.list.map((item) => {
			item.date = new Date(item.date).toISOString().substring(0, 10);
			item.changes = item.changes.sort((a, b) => a.file.localeCompare(b.file));
			return item;
		});
		total.value = list.total;

		state.value = LoadState.Loaded;

		return updateList;
	} catch (error) {
		state.value = LoadState.Error;
	}

	return [];
}

export function useUpdates() {
	return {
		state,
		total,
		getUpdates,
	};
}
