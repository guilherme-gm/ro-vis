import { ParsingResult } from "./ParsingResult.js";

export type ParserResult<T> = {
	result: ParsingResult;
	data: T[];
};

export interface IParser<T> {
	parse(): Promise<ParserResult<T>>;
}
