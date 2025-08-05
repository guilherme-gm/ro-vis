import type { MapNpc } from "@/models/Map";
import type { MapNpcDiffer } from "./differs";

export type DiffRecord = {
	id: string;
	diffType: 'added' | 'removed' | 'unchanged';
	value: string;
};

export type DiffedList = DiffRecord[];

export class MapNpcFormatter {
	public createId = (npc: MapNpc) => {
		const x = npc.Location.X.toString().padStart(3, '0');
		const y = npc.Location.Y.toString().padStart(3, '0');

		return `${x}#${y}#${npc.Type}#${npc.SpriteId}#${npc.Name1.Value}#${npc.Name2}`;
	}

	public formatList = (npcs?: MapNpc[]): string[] => {
		if (!npcs) {
			return [];
		}

		return npcs
			.sort((a, b) => this.createId(a).localeCompare(this.createId(b)))
			.map((npc) => this.format(npc));
	};

	public format = (npc: MapNpc): string => {
		return `[${npc.Location.X}, ${npc.Location.Y}] ${npc.Name1.Value} (${npc.Name2})`;
	};

	static use(): MapNpcFormatter {
		return new MapNpcFormatter();
	}
}

export class MapNpcDiffFormatter {
	private differ: MapNpcDiffer;
	private formatter: MapNpcFormatter;

	constructor(differ: MapNpcDiffer) {
		this.differ = differ;
		this.formatter = MapNpcFormatter.use();
	}

	public formatList = (): DiffedList => {
		const diffedList: DiffedList = [];

		const addedEntries = this.differ.diff.value.added
			.map((npc): DiffRecord => ({
				id: this.formatter.createId(npc),
				value: this.formatter.format(npc),
				diffType: 'added'
			}));

		const removedEntries = this.differ.diff.value.removed
			.map((npc): DiffRecord => ({
				id: this.formatter.createId(npc),
				value: this.formatter.format(npc),
				diffType: 'removed'
			}));

		const unchangedEntries = this.differ.diff.value.unchanged
			.map((npc): DiffRecord => ({
				id: this.formatter.createId(npc),
				value: this.formatter.format(npc),
				diffType: 'unchanged'
			}));

		diffedList.push(...addedEntries)
		diffedList.push(...removedEntries)
		diffedList.push(...unchangedEntries)

		return diffedList.sort((a, b) => a.id.localeCompare(b.id));
	};

	static use(differ: MapNpcDiffer): MapNpcDiffFormatter {
		return new MapNpcDiffFormatter(differ);
	}
}
