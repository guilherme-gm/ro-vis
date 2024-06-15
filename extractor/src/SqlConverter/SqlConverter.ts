import * as fs from "fs/promises";
import { SqlFieldMetadata } from "./Decorators/SqlFieldMetadata.js";
import path from "path";

export class SqlConverter {
	protected valueToSql(value: unknown, type: unknown): string {
		if (type === Number) {
			return `${value}`;
		}

		if (type === String) {
			return `"${(value as string).replace(/"/g, '""').replace(/\n/g, '\\n')}"`;
		}

		if (type === Boolean) {
			return value ? 'TRUE' : 'FALSE';
		}

		throw new Error(`Can't convert value ${value} to type ${type}.`);
	}

	protected convertInto(data: object, converted: [string, string][], prefix: string = ''): void {
		const fields = Reflect.getMetadata('sql-field', data) as SqlFieldMetadata[] | undefined;
		if (!fields) {
			return;
		}

		fields.forEach((field) => {
			if (field.type === 'simple') {
				const fieldType = field.outType?.() ?? Reflect.getMetadata('design:type', data, field.key);
				const fieldName = field.key.toString();

				// @ts-expect-error -- too hard to type
				let value = data[field.key];
				if (field.transform) {
					value = field.transform(value);
				}

				const sqlFieldName = prefix ? `${prefix}_${fieldName}` : fieldName;

				converted.push([sqlFieldName, this.valueToSql(value, fieldType)]);
			} else {
				const newPrefix = prefix ? `${prefix}_${field.key.toString()}` : field.key;
				// @ts-expect-error -- too hard to type
				this.convertInto(data[field.key], converted, newPrefix);
			}
		});
	}

	protected generateReplace(tableName: string, fields: [string, string][]): string {
		const fieldNames = fields
			.map((field) => `\`${field[0]}\``)
			.join(', ');
		const fieldValues = fields
			.map((field) => `${field[1]}`)
			.join(', ');

		return `REPLACE INTO \`${tableName}\` (${fieldNames}) VALUES (${fieldValues});`;
	}

	public async convert(table: string, data: object[]): Promise<void> {
		const replaces: string[] = [];

		data.forEach((row) => {
			const convertedPairs: [string, string][] = [];
			this.convertInto(row, convertedPairs);

			const sql = this.generateReplace(table, convertedPairs);
			replaces.push(sql);
		});

		await fs.mkdir(path.resolve("out", "sql"), { recursive: true });
		await fs.writeFile(path.resolve("out", "sql", `${table}.sql`), replaces.join('\n'));
	}
}
