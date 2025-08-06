import { computed, type ComputedRef } from "vue";
import type { Formatter } from "./formatters";

type DiffResult<T> = {
	added: T[];
	removed: T[];
	unchanged: T[];
}

export class ListDiffer<T> {
	private from: ComputedRef<T[]>;

	private to: ComputedRef<T[]>;

	private formatter: Formatter<T>;

	constructor(from: ComputedRef<T[]>, to: ComputedRef<T[]>, formatter: Formatter<T>) {
		this.from = from;
		this.to = to;
		this.formatter = formatter;
	}

	public diff = computed((): DiffResult<T> => {
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

	static use<T>(from: ComputedRef<T[]>, to: ComputedRef<T[]>, formatter: Formatter<T>): ListDiffer<T> {
		return new ListDiffer(from, to, formatter);
	}
}
