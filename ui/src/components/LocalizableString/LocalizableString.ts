import type { LocalizableString } from "@/models/LocalizableString";
import { computed, ref, type Ref } from "vue";
import { StringLoader } from "./StringLoader";

/**
 * Handles LocalizableString content loading.
 * It leverages StringLoader for that
 */
export class LocalizableStringComponent {
	private val: Ref<LocalizableString>;

	private status = ref<'loading' | 'loaded' | 'error'>('loading');

	private stringLoader = StringLoader.use();

	public constructor(val: Ref<LocalizableString>) {
		this.val = val;

		if (this.val.value.Kind === 'localized') {
			this.load(this.val.value.Value);
		} else {
			console.error('StringLoader: String is not localized');
		}
	}

	private load = async (id: string) => {
		try {
			await this.stringLoader.enqueueLoad(id);
			this.status.value = 'loaded';
		} catch (error) {
			console.error(error);
			this.status.value = 'error';
		}
	}

	public displayString = computed(() => {
		switch (this.status.value) {
			case 'loaded':
				return this.stringLoader.getString(this.val.value.Value);
			case 'loading':
				return `[Loading... #${this.val.value.Value}]`;
			case 'error':
				return `[Error loading #${this.val.value.Value}]`;
			default:
				console.error('StringLoader: Unknown status');
				return this.val.value.Value;
		}
	});

	public static use(val: Ref<LocalizableString>): LocalizableStringComponent {
		return new LocalizableStringComponent(val);
	}
}
