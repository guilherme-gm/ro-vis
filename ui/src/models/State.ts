export enum StatePriority {
	None = -1,
	Gold = 0,
	Red = 1,
	Blue = 2,
	Green = 3,
	White = 4,
}

export type DescriptionLine = {
	Text: string;

	Color: number[] | null;
}

/**
 * Represents a State record in the tool.
 */
export type State = {
	Patch: string;

	/**
	 * The File Version that originated this object
	 */
	_FileVersion: number;

	/**
	 * State ID (SI / EFST_)
	 */
	Id: number;

	/**
	 * EFST_ constant
	 */
	Constant: string;

	Description: DescriptionLine[];

	HasTimeLimit: boolean;

	/**
	 * The number of the line in Description containing the time text.
	 * NOTE: This is 1-based, thus TimeStrLineNum = 1 means Description[0]
	 *       It is "-1" if there is no string in this entry.
	 */
	TimeStrLineNum: number;

	HasEffectImage: boolean;

	IconImage: string;

	IconPriority: StatePriority;
}

export type MinState = {
	Id: string;

	Patch: string;

	Constant: string;
};
