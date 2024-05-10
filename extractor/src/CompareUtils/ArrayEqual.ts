export class ArrayEqual {
	public static isEqual<T>(array1: T[], array2: T[], fn: (a: T, b: T) => boolean = (a, b) => a === b): boolean {
		if (array1.length !== array2.length) {
			return false;
		}

		for (let i = 0; i < array1.length; i++) {
			if (!fn(array1[i]!, array2[i]!)) {
				return false;
			}
		}

		return true;
	}
}
