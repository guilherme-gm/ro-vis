import { CommonApi, type PatchItem } from "./CommonApi";
import type { Map, MinMap } from "@/models/Map";

export type MapPatch = PatchItem<Map>;

export class MapApi extends CommonApi<Map, MinMap> {
	public static use() {
		return new MapApi('maps/');
	}
}
