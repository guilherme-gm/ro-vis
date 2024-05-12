import { ClassConstructor, instanceToPlain, plainToInstance } from "class-transformer";

export class ConvertClass {
	public static convert<T>(value: unknown, to: ClassConstructor<T>): T {
		return plainToInstance(to, instanceToPlain(value), {
			excludeExtraneousValues: true,
			exposeDefaultValues: true,
		});
	}
}
