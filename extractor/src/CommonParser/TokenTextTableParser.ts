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

	private isEscapedSeparator(): boolean {
		// In multiline cells, a "#" may be followed by a space and other characters, in this case they are considered part
		// of the text and not a delimiter.
		return this.reader.peekChar() === '#'
			&& this.reader.peekChar(1) === ' '
			&& /[\w\d]+/.test(this.reader.peekChar(2) ?? '');
	}

	private _readCell(multiline: boolean): Result<string> {
		this.skipIgnorable();

		if (this.isEndOfFile()) {
			return Result.fail(new Error('End of file'));
		}

		let cell = '';
		while (!this.isEndOfFile()
			&& (multiline || this.reader.peekChar() !== '\n')
			&& (this.reader.peekChar() !== '#' || this.isEscapedSeparator())
		) {
			cell += this.reader.readChar();
		}

		if (this.isEndOfFile()) {
			return Result.fail(new Error(`End of file reached before finding the end of a cell. Value: "${cell}". Line ${this.reader.getLineNumber()}`));
		}

		const char = this.reader.readChar();

		if (char === '\n' && !multiline) {
			if (cell.trim() === '') {
				this.reader.readChar();
				return this._readCell(multiline);
			}

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

		const text = cell.unwrap();
		const value = parseInt(text, 10);
		if (isNaN(value)) {
			throw new Error(`Detected an invalid number in a cell. Value: "${text}". Line ${this.reader.getLineNumber()}`);
		}

		return Result.ok(value);
	}

	public readMultilineCell(): Result<string> {
		const result = this._readCell(true);
		if (!result.isOk()) {
			return result;
		}

		const text = result.unwrap().trim();
		return Result.ok(text);
	}
}
