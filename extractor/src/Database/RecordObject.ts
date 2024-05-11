export interface RecordObject {
	getId(): string;
	getFileVersion(): number;
	equals(other: RecordObject): boolean;
}
