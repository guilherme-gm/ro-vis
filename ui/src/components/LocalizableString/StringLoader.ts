import { I18nApi } from "@/services/I18nApi";

/**
 * The actual process to load LocalizableStrings.
 * It does that by enqueueing string requests, and submitting it after a few miliseconds,
 * this enables loading many strings at once, and also keeping a local cache for the session.
 *
 * This has a few advantages:
 * - Loading many strings at once is faster and reduces the amount of requests to the server
 * - Keeping a local cache for the session is faster than loading them from the server every time
 *
 * On the other hand, the code becomes a little more complex and has a small delay.
 */
export class StringLoader {
	private static instance: StringLoader | null = null;

	private strings = new Map<string, string>();

	private queue = new Set<string>();

	private timer: number | null = null;

	private queueResolution: Promise<void> | null = null;


	public static use(): StringLoader {
		if (!StringLoader.instance) {
			StringLoader.instance = new StringLoader();
		}

		return StringLoader.instance;
	}

	private loadQueue = async (ids: Set<string>) => {
		const strings = await I18nApi.use().getStrings([...ids]);
		for (const { I18nId, PtBrText } of strings) {
			this.strings.set(I18nId, PtBrText);
		}
	}

	public getString = (id: string): string => {
		return this.strings.get(id) ?? id;
	}

	public enqueueLoad = async (id: string) => {
		if (this.strings.has(id)) {
			return;
		}

		this.queue.add(id);

		if (this.timer) {
			await this.queueResolution;
			return;
		}

		this.queueResolution = new Promise((resolve, reject) => {
			this.timer = setTimeout(async () => {
				const queue = this.queue;
				this.queue = new Set<string>();
				this.timer = null;
				this.queueResolution = null;

				try {
					await this.loadQueue(queue);
				} catch (error) {
					console.error(error);
					reject(error);
					return;
				}
				resolve();
			}, 10);
		});
	}
}
