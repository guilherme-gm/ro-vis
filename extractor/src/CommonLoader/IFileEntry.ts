export interface IFileEntry<T> {
	getId(): string;

	toEntity(): T;
}
