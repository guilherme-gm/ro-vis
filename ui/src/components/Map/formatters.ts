import type { MapNpc } from "@/models/Map";
import type { ListDiffer } from "./differs";

export type DiffRecord = {
	id: string;
	diffType: 'added' | 'removed' | 'unchanged';
	value: string;
};

export type DiffedList = DiffRecord[];

export interface Formatter<T> {
	createId(item: T): string;
	format(item: T): string;
}

export class MapNpcFormatter implements Formatter<MapNpc> {
	public createId = (npc: MapNpc) => {
		const x = npc.Location.X.toString().padStart(3, '0');
		const y = npc.Location.Y.toString().padStart(3, '0');

		return `${x}#${y}#${npc.Type}#${npc.SpriteId}#${npc.Name1.Value}#${npc.Name2}`;
	}

	public format = (npc: MapNpc): string => {
		return `[${npc.Location.X}, ${npc.Location.Y}] ${npc.Name1.Value} (${npc.Name2})`;
	};

	static use(): MapNpcFormatter {
		return new MapNpcFormatter();
	}
}

export class ListFormatter<T> {
	private formatter: Formatter<T>;

	constructor(formatter: Formatter<T>) {
		this.formatter = formatter;
	}

	public formatList = (items?: T[]): string[] => {
		if (!items) {
			return [];
		}

		return items
			.sort((a, b) => this.formatter.createId(a).localeCompare(this.formatter.createId(b)))
			.map((item) => this.formatter.format(item));
	};

	static use<T>(formatter: Formatter<T>): ListFormatter<T> {
		return new ListFormatter<T>(formatter);
	}
}

export class ListDiffFormatter<T> {
	private differ: ListDiffer<T>;
	private formatter: ListFormatter<T>;
	private specificFormatter: Formatter<T>;

	constructor(differ: ListDiffer<T>, formatter: ListFormatter<T>, specificFormatter: Formatter<T>) {
		this.differ = differ;
		this.formatter = formatter;
		this.specificFormatter = specificFormatter;
	}

	public formatList = (): DiffedList => {
		const diffedList: DiffedList = [];

		const addedEntries = this.differ.diff.value.added
			.map((npc): DiffRecord => ({
				id: this.specificFormatter.createId(npc),
				value: this.specificFormatter.format(npc),
				diffType: 'added'
			}));

		const removedEntries = this.differ.diff.value.removed
			.map((npc): DiffRecord => ({
				id: this.specificFormatter.createId(npc),
				value: this.specificFormatter.format(npc),
				diffType: 'removed'
			}));

		const unchangedEntries = this.differ.diff.value.unchanged
			.map((npc): DiffRecord => ({
				id: this.specificFormatter.createId(npc),
				value: this.specificFormatter.format(npc),
				diffType: 'unchanged'
			}));

		diffedList.push(...addedEntries)
		diffedList.push(...removedEntries)
		diffedList.push(...unchangedEntries)

		return diffedList.sort((a, b) => a.id.localeCompare(b.id));
	};

	static use<T>(differ: ListDiffer<T>, formatter: ListFormatter<T>, specificFormatter: Formatter<T>): ListDiffFormatter<T> {
		return new ListDiffFormatter(differ, formatter, specificFormatter);
	}
}
