import chalk from 'chalk';

export class Logger {
	public static logger = new Logger();

	public static status(message: string): void {
		Logger.logger.status(message);
	}

	public static info(message: string): void {
		Logger.logger.info(message);
	}

	public static error(message: string, error?: unknown): void {
		Logger.logger.error(message, error);
	}

	public static warn(message: string): void {
		Logger.logger.warn(message);
	}

	public status(message: string): void {
		console.log(chalk.dim(`${chalk.green('[Status]')} ${message}`));
	}

	public info(message: string): void {
		console.log(chalk.dim(`${chalk.whiteBright('[Info]')} ${message}`));
	}

	public error(message: string, error?: unknown): void {
		console.error(chalk.dim(`${chalk.red('[Error]')} ${message}`), error ?? '');
	}

	public warn(message: string): void {
		console.warn(chalk.dim(`${chalk.yellow('[Warn]')} ${message}`));
	}
}
