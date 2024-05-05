export class Result<T> {
	constructor(
		public readonly value: T | undefined,
		public readonly error: Error | undefined
	) { }

	public static ok<T>(value: T): Result<T> {
		return new Result<T>(value, undefined);
	}

	public static fail<T>(error: Error): Result<T> {
		return new Result<T>(undefined, error);
	}

	public isOk(): boolean {
		return this.error === undefined;
	}

	public isFail(): boolean {
		return this.error !== undefined;
	}

	public unwrap(): T {
		if (this.isFail()) {
			throw this.error;
		}

		return this.value!;
	}
}
