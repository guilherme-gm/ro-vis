import { CommonApi, type PatchItem } from "./CommonApi";
import type { I18n, MinI18n } from "@/models/I18n";

export type I18nPatch = PatchItem<I18n>;

export class I18nApi extends CommonApi<I18n, MinI18n> {
	public getStrings = async (ids: string[]) => {
		const res = await this.api.post<{ strings: { I18nId: string; PtBrText: string }[] }>(this.path + 'text/', { ids });
		console.log(res);
		return res.strings;
	}

	public static use() {
		return new I18nApi('i18n/');
	}
}
