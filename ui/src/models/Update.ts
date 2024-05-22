import type { Patch } from "./Patch";

export type UpdateItem = {
	file: string;
	patch: string;
};

export type Update = {
	id: string;

	order: number

	updates: UpdateItem[];

	patches: Patch[] | null;

	tags: string[];
}
