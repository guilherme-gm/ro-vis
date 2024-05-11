export class CharStreamReader {
	private readonly text: string = '';

	private lineNumber: number = 1;

	private pointer: number = 0;

	constructor (text: string) {
		this.text = text.replace(/\r/g, "").replace(/\u3000/g, "");
	}

	public isEndOfFile(): boolean {
		return this.pointer >= this.text.length;
	}

	public isLineStart(): boolean {
		if (this.pointer == 0) {
			return true;
		}

		if (this.text[this.pointer - 1] === '\n') {
			return true;
		}

		return false;
	}

	public getLineNumber(): number {
		return this.lineNumber;
	}

	public readChar(): string {
		if (this.isEndOfFile()) {
			throw new Error('End of file reached.');
		}

		const char = this.text[this.pointer]!;
		this.pointer++;

		if (char === '\n') {
			this.lineNumber++;
		}

		return char;
	}

	public peekChar(skip: number = 0): string | null {
		if (this.isEndOfFile()) {
			return null;
		}

		if (this.text.length <= this.pointer + skip) {
			return null;
		}

		return this.text[this.pointer + skip]!;
	}

	public unreadChar(): void {
		if (this.pointer == 0) {
			throw new Error('Cannot unread. Already at the beggining of the file.');
		}

		this.pointer--;
		const char = this.text[this.pointer];

		if (char === '\n') {
			this.lineNumber--;
		}
	}
}
