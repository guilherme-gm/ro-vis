import { SqlFieldMetadata } from "./SqlFieldMetadata.js";

export function SqlField(params?: Pick<SqlFieldMetadata, 'transform' | 'outType' | 'name'>): PropertyDecorator {
	return (target, key) => {
		const fields = (Reflect.getMetadata('sql-field', target) || []) as SqlFieldMetadata[];

		const exists = fields.some((field) => field.key === key);
		if (!exists) {
			fields.push({
				key,
				type: 'simple',
				...params
			});
		}
		Reflect.defineMetadata('sql-field', fields, target)
	}
}
