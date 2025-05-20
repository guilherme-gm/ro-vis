export type UpdateChange = {
	file: string;
	patch: string;
};

export type Update = {
	date: string;

	changes: UpdateChange[];
}
