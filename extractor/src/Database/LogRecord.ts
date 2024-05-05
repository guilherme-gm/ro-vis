import { RecordObject } from './RecordObject.js';

export class LogRecord<T extends RecordObject> {
	public _id!: string;

	public current!: {
		patch: string;
		value: T;
	};

	public previous!: {
		patch: string;
		value: T;
	}[];

	constructor(patch: string, currentValue: T);
	constructor(record: LogRecord<T>);
	constructor(param1: string | LogRecord<T>, currentValue?: T) {
		if (currentValue) {
			const patch = param1 as string;
			this._id = currentValue.getId();
			this.current = {
				patch,
				value: currentValue,
			};
			this.previous = [];
		} else {
			const record = param1 as LogRecord<T>;
			this._id = record._id;
			this.current = record.current;
			this.previous = record.previous;
		}
	}

	public addChange(patch: string, value: T): void {
		if (this._id !== value.getId()) {
			throw new Error(`ID Mismatch (${this._id} vs ${value.getId()})`);
		}

		if (this.current.patch === patch) {
			throw new Error(`${this._id} already has a change for patch "${patch}.`);
		}

		this.previous.push(this.current);
		this.current = {
			patch,
			value,
		};
	}
}
