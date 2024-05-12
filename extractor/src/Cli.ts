import { Command } from "commander";

export class Cli {
	public static cli: Cli = new Cli();

	public static load(): void {
		Cli.cli.load();
	}

	private program = new Command();

	public dryRun: boolean = false;

	public cleanRun: boolean = false;

	public mongoPort: number | undefined;

	public holdProcess: boolean = false;

	public only: string = "";

	public load(): void {
		this.program.name('RO Vis extractor');

		this.program
			.option("--dry-run", "Replicate data to an in-memory DB to execute the process and dump to file")
			.option("--clean-run", "To be used together with '--dry-run'. Only replicate patch list, everything else is run from zero")
			.option("--mongo-port <port>", "For dry-run and clean-run, defines the port where the in-memory MongoDB will run on. When not set, a random one is chosen.")
			.option("--hold-process", "When the extraction finishes/crashes, stop in a 'Press enter to continue' message before ending the process (the temporary DB is also kept running)")
			.option("--only <name>", "Run only a given loader");

		this.program.parse(process.argv);

		this.dryRun = this.program.opts()['dryRun'] ?? false;
		this.cleanRun = this.program.opts()['cleanRun'] ?? false;
		this.mongoPort = this.program.opts()['mongoPort'] ? parseInt(this.program.opts()['mongoPort'], 10) : undefined;
		this.holdProcess = this.program.opts()['holdProcess'] ?? false;
		this.only = this.program.opts()['only'] ?? '';
	}
}