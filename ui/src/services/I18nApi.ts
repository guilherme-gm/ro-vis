import { CommonApi, type PatchItem } from "./CommonApi";
import type { I18n, MinI18n } from "@/models/I18n";

export type I18nPatch = PatchItem<I18n>;

export class I18nApi extends CommonApi<I18n, MinI18n> {
	public static use() {
		return new I18nApi('i18n/');
	}
}
