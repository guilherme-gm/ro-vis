export type SqlFieldMetadata = {
	key: string | symbol;
	type: 'simple' | 'nested';
	name?: string | undefined;
	transform?: ((value: any) => string) | undefined;
	outType?: (() => unknown) | undefined;
}
