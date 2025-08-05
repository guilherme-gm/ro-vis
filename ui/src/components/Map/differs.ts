import type { MapNpc } from "@/models/Map";
import { computed, type ComputedRef } from "vue";
import { MapNpcFormatter } from "./formatters";

type DiffResult<T> = {
	added: T[];
	removed: T[];
	unchanged: T[];
}

export class MapNpcDiffer {
	private from: ComputedRef<MapNpc[]>;

	private to: ComputedRef<MapNpc[]>;

	private formatter: MapNpcFormatter;

	constructor(from: ComputedRef<MapNpc[]>, to: ComputedRef<MapNpc[]>) {
		this.from = from;
		this.to = to;
		this.formatter = MapNpcFormatter.use();
	}

	public diff = computed((): DiffResult<MapNpc> => {
		if (!this.from.value) {
			return { added: this.to.value ?? [], removed: [], unchanged: [] };
		}

		if (!this.to.value) {
			return { added: [], removed: this.from.value ?? [], unchanged: [] };
		}

		const existing = this.from.value;
		const existingIds = new Set(existing.map((r) => this.formatter.createId(r)));
		const newIds = new Set(this.to.value.map((r) => this.formatter.createId(r)));

		const added = this.to.value.filter((r) => !existingIds.has(this.formatter.createId(r)));
		const removed = this.from.value.filter((r) => !newIds.has(this.formatter.createId(r)));
		const unchanged = this.from.value.filter((r) => existingIds.has(this.formatter.createId(r)));

		return { added, removed, unchanged };
	});

	static use(from: ComputedRef<MapNpc[]>, to: ComputedRef<MapNpc[]>): MapNpcDiffer {
		return new MapNpcDiffer(from, to);
	}
}
