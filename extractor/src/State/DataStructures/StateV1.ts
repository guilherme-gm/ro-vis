import { Expose } from "class-transformer";
import { IFileEntry } from "../../CommonLoader/IFileEntry.js";
import { StateV } from "./StateV.js";
import { State } from "./State.js";
import { DescriptionLine } from "./DescriptionLine.js";
import { ConvertClass } from "../../Utils/ConvertClass.js";
import { StatePriority } from "./StatePriority.js";

/**
 * efstids.lub, stateiconinfo_f.lub, stateiconimginfo.lub, stateiconinfo.lub
 * Since before 2012.
 */
export class StateV1 implements IFileEntry<State> {
	public static isV1(quest: StateV): quest is StateV1 {
		return quest._FileVersion === 1;
	}

	@Expose()
	public readonly _FileVersion: number = 1;

	/**
	 * State ID (SI / EFST_)
	 */
	@Expose()
	public Id: number = 0;

	/**
	 * EFST_ constant
	 */
	@Expose()
	public Constant: string = "";

	@Expose()
	public Description: DescriptionLine[] = [];

	@Expose()
	public HasTimeLimit: boolean = false;

	@Expose()
	public TimeStrLineNum: number = -1;

	/**
	 * Whether it has an effect image (HaveEfstImgTable)
	 */
	@Expose()
	public HasEffectImage: boolean = false;

	@Expose()
	public IconImage: string = "";

	@Expose()
	public IconPriority: StatePriority = StatePriority.None;

	public static fromState(state: State): StateV1 {
		if (state._FileVersion > 1) {
			throw new Error(`Can not convert item v${state._FileVersion} to V1`);
		}

		return ConvertClass.convert(state, StateV1);
	}

	public getId(): string {
		return this.Id.toString();
	}

	public getFileVersion(): number {
		return this._FileVersion;
	}

	public toEntity(): State {
		return ConvertClass.convert(this, State);
	}
}
