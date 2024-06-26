import { CommonApi, type PatchItem } from "./CommonApi";
import type { MinState, State } from "@/models/State";

export type StatePatch = PatchItem<State>;

export class StateApi extends CommonApi<State, MinState> {
	public static use() {
		return new StateApi('states/');
	}
}
