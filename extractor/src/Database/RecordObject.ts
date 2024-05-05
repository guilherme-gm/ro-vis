export interface RecordObject {
	getId(): string;
	getFileVersion(): number;
	hasChange(other: RecordObject): boolean;
}
