import chalk from 'chalk';
import * as fs from "fs/promises";
import iconv from "iconv-lite";
import { Logger } from '../../../Logger.js';
import { ItemMoveInfoV5 } from '../../DataStructures/ItemMoveInfoV5.js';

export class ItemMoveInfoV5Parser {
	public static async fromFile(filePath: string): Promise<ItemMoveInfoV5Parser> {
		const rawContent = await fs.readFile(filePath);
		const content = iconv.decode(rawContent, 'euc-kr').toString();
		return new ItemMoveInfoV5Parser(content);
	}

	private lines: string[];

	protected constructor(content: string) {
		this.lines = content.replace(/\r\n/g, "\n").split('\n');
	}

	private parseRestriction(column: string): boolean {
		// 0 = can do it ; 1 = can't do it
		if (column !== '0' && column !== '1') {
			throw new Error(`Unexpected ItemMoveInfo value: "${column}".`);
		}

		return column === '0';
	}

	public parse(): Promise<ItemMoveInfoV5[]> {
		const table: ItemMoveInfoV5[] = [];

		this.lines.forEach((line, idx) => {
			line = line.trim();
			if (line.length === 0 || line.startsWith("//")) {
				return;
			}

			const columns = line.split('\t').map((v) => v.trim());
			if (columns.length !== 10) {
				Logger.warn(`${chalk.whiteBright("ItemMove:")} Line ${idx + 1} has ${columns.length} columns, 10 expected. Skipping...`);
				return;
			}

			const moveInfo = new ItemMoveInfoV5();
			moveInfo.itemId = parseInt(columns[0]!, 10);
			moveInfo.canDrop = this.parseRestriction(columns[1]!);
			moveInfo.canTrade = this.parseRestriction(columns[2]!);
			moveInfo.canMoveToStorage = this.parseRestriction(columns[3]!);
			moveInfo.canMoveToCart = this.parseRestriction(columns[4]!);
			moveInfo.canSellToNpc = this.parseRestriction(columns[5]!);
			moveInfo.canMail = this.parseRestriction(columns[6]!);
			moveInfo.canAuction = this.parseRestriction(columns[7]!);
			moveInfo.canMoveToGuildStorage = this.parseRestriction(columns[8]!);
			moveInfo.commentName = columns[9]?.replace(/^\/\/ /, "") ?? "";
		});

		return Promise.resolve(table);
	}
}
