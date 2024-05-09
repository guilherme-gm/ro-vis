export interface IParser<T> {
	parse(): Promise<T[]>;
}
