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
		console.log(`${chalk.green('[Status]')} ${message}`);
	}

	public info(message: string): void {
		console.log(`${chalk.whiteBright('[Info]')} ${message}`);
	}

	public error(message: string, error?: unknown): void {
		if (error) {
			console.error(`${chalk.red('[Error]')} ----------- ${message} -----------`);
			console.error(error);
			console.error(`${chalk.red('[Error]')} -------------------------------------------`);
			console.error('');
		} else {
			console.error(`${chalk.red('[Error]')} ${message}`);
		}
	}

	public warn(message: string): void {
		console.warn(`${chalk.yellow('[Warn]')} ${message}`);
	}
}
