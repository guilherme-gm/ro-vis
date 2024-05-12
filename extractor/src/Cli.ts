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

	public changeDump: boolean = false;

	public only: string = "";

	public onlyPatches: string[] | null = [];

	public showTraces: boolean = false;

	public load(): void {
		this.program.name('RO Vis extractor');

		this.program
			.option("--dry-run", "Replicate data to an in-memory DB to execute the process")
			.option("--clean-run", "Run everything from scratch in a in-memory DB")
			.option("--change-dump", "Dump change logs as json files for debugging")
			.option("--mongo-port <port>", "For dry-run and clean-run, defines the port where the in-memory MongoDB will run on. When not set, a random one is chosen.")
			.option("--hold-process", "When the extraction finishes/crashes, stop in a 'Press enter to continue' message before ending the process (the temporary DB is also kept running)")
			.option("--only <name>", "Run only a given loader")
			.option("--only-patches <patches>", "Only run the following patches, several patches may be specified separating with comma")
			.option("--show-traces", "Include error traces? When not passed, only the error message is shown.");

		this.program.parse(process.argv);

		this.dryRun = this.program.opts()['dryRun'] ?? false;
		this.cleanRun = this.program.opts()['cleanRun'] ?? false;
		this.changeDump = this.program.opts()['changeDump'] ?? false;
		this.mongoPort = this.program.opts()['mongoPort'] ? parseInt(this.program.opts()['mongoPort'], 10) : undefined;
		this.holdProcess = this.program.opts()['holdProcess'] ?? false;
		this.only = this.program.opts()['only'] ?? '';
		this.onlyPatches = this.program.opts()['onlyPatches']?.split(',') ?? null;
		this.showTraces = this.program.opts()['showTraces'] ?? false;
	}
}
