import { ref } from "vue";
import { useApi } from "./api";
import { LoadState } from "./LoadState";

export type Record<T> = {
	Update: string;
	Data: T;
}

export type PatchItem<T> = {
	LastUpdated: string | null;
	From: Record<T> | null;
	To: Record<T> | null;
};

export type PaginatedResponse<T> = {
	total: number;
	list: T[];
};

export class CommonApi<Entity, MinEntity> {
	public total = ref(-1);

	public historyTotal = ref(-1);

	public state = ref<LoadState>(LoadState.None);

	private api = useApi();

	private path: string;

	constructor(path: string) {
		this.path = path;
		if (!this.path.endsWith('/')) {
			this.path += '/';
		}
	}

	public getItems = async (page: number): Promise<PaginatedResponse<MinEntity>['list']> => {
		try {
			this.state.value = LoadState.Loading;

			const items = await this.api.get<PaginatedResponse<MinEntity>>(this.path, { start: ((page - 1) * 100) });

			this.total.value = items.total;

			this.state.value = LoadState.Loaded;

			return items.list;
		} catch (error) {
			this.state.value = LoadState.Error;
		}

		return [];
	};

	public getItemHistory = async (id: string, page: number): Promise<PatchItem<Entity>[]> => {
		try {
			this.state.value = LoadState.Loading;

			const items = await this.api.get<PaginatedResponse<PatchItem<Entity>>>(`${this.path}/${id}`, { start: ((page - 1) * 100) });

			this.historyTotal.value = items.total;

			this.state.value = LoadState.Loaded;

			return items.list;
		} catch (error) {
			this.state.value = LoadState.Error;
		}

		return [];
	};

	public getPatchItems = async (patch: string, page: number): Promise<PatchItem<Entity>[]> => {
		try {
			this.state.value = LoadState.Loading;

			const items = await this.api.get<PaginatedResponse<PatchItem<Entity>>>(`${this.path}update/${patch}`, { start: ((page - 1) * 100) });

			this.total.value = items.total;

			this.state.value = LoadState.Loaded;

			return items.list;
		} catch (error) {
			this.state.value = LoadState.Error;
		}

		return [];
	};

	public static use(path: string) {
		return new CommonApi(path);
	}
}
