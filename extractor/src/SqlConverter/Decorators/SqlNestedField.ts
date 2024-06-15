import { SqlFieldMetadata } from "./SqlFieldMetadata.js";

export function SqlNestedField(params?: { transform?: SqlFieldMetadata['transform']; outType?: SqlFieldMetadata['outType']; }): PropertyDecorator {
	return (target, key) => {
		const fields = (Reflect.getMetadata('sql-field', target) || []) as SqlFieldMetadata[];

		const exists = fields.some((field) => field.key === key);
		if (!exists) {
			fields.push({ key, type: 'nested', transform: params?.transform, outType: params?.outType });
		}
		Reflect.defineMetadata('sql-field', fields, target)
	}
}
