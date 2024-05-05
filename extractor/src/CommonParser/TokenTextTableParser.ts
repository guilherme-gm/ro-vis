import { CharStreamReader } from "./CharStreamReader.js";
import { Result } from "./Result.js";

export class TokenTextTableParser {
	private readonly reader: CharStreamReader;

	protected constructor(text: string) {
		this.reader = new CharStreamReader(text);
	}

	private skipComment(): boolean {
		if (!this.reader.isLineStart()) {
			return false;
		}

		if (this.reader.peekChar() !== '/') {
			return false;
		}

		while (!this.reader.isEndOfFile() && this.reader.readChar() !== '\n') {
			continue;
		}

		return true;
	}

	private skipIgnorable(): void {
		let stop = false;
		while (!this.reader.isEndOfFile() && !stop) {
			if (this.skipComment()) {
				continue;
			}

			const char = this.reader.readChar();
			stop = char !== '\n' && char !== ' ';
			if (stop) {
				this.reader.unreadChar();
			}
		}
	}

	public isEndOfFile(): boolean {
		return this.reader.isEndOfFile();
	}

	public isLineStart(): boolean {
		return this.reader.isLineStart();
	}

	public consumeRestOfLine(): void {
		while (!this.reader.isEndOfFile() && this.reader.readChar() !== '\n') {
			continue;
		}
	}

	private _readCell(multiline: boolean): Result<string> {
		this.skipIgnorable();

		if (this.isEndOfFile()) {
			return Result.fail(new Error('End of file'));
		}

		let cell = '';
		while (!this.isEndOfFile()
			&& (multiline || this.reader.peekChar() !== '\n')
			&& this.reader.peekChar() !== '#'
		) {
			cell += this.reader.readChar();
		}

		if (this.isEndOfFile()) {
			return Result.fail(new Error(`End of file reached before finding the end of a cell. Value: "${cell}". Line ${this.reader.getLineNumber()}`));
		}

		const char = this.reader.readChar();

		if (char === '\n' && !multiline) {
			return Result.fail(new Error(`Detected a line break in a cell that does not support multi line. Value: "${cell}". Line ${this.reader.getLineNumber()}`));
		}

		this.skipIgnorable();

		return Result.ok(cell);
	}

	public readCell(): Result<string> {
		return this._readCell(false);
	}

	public readIntCell(): Result<number> {
		const cell = this._readCell(false);
		if (cell.isFail()) {
			return Result.fail(cell.error!);
		}

		const value = parseInt(cell.unwrap(), 10);
		if (isNaN(value)) {
			throw new Error(`Detected an invalid number in a cell. Value: "${cell}". Line ${this.reader.getLineNumber()}`);
		}

		return Result.ok(value);
	}

	public readMultilineCell(): Result<string> {
		return this._readCell(true);
	}
}
