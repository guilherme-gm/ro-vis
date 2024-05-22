import { ref } from "vue";
import type { Update } from "@/models/Update";
import { useApi } from "./api";

const api = useApi();

const isLoading = ref(false);
const updateList = ref<Update[]>([]);

async function getUpdates(): Promise<void> {
	try {
		isLoading.value = true;

		const list = await api.get<Update[]>('updates/');

		updateList.value = list.map((item) => {
			item.updates = item.updates.sort((a, b) => a.file.localeCompare(b.file));
			return item;
		});
	} catch (error) {
		console.error(error);
	} finally {
		isLoading.value = false;
	}
}

export function useUpdates() {
	return {
		isLoading,
		updateList,
		getUpdates,
	};
}
