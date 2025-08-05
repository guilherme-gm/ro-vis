import type { LocalizableString } from "./LocalizableString";

export type Map = {
	HistoryID: string;

	FileVersion: number;

	Id: string;

	Name: string;

	SpecialCode: number;

	Mp3Name: string;

	Npcs: MapNpc[];

	Warps: MapWarp[];

	Spawns: MapSpawn[];
}

export type MapNpc = {
	Type: string;

	SpriteId: number;

	Name1: LocalizableString;

	Name2: string;

	Location: MapCoord;
}

export type MapCoord = {
	MapId: string;

	X: number;

	Y: number;
}

export type MapWarp = {
	From: MapCoord;

	To: MapCoord;

	WarpType: string;

	SpriteId: number;

	Name1: LocalizableString;

	Name2: string;
}

export type MapSpawn = {
	Type: string;

	SpriteId: number;

	Name1: LocalizableString;

	Name2: string;

	Level: number;

	Amount: number;

	Element: number;

	Size: number;

	Race: number;
}

export type MinMap = {
	MapID: string;

	LastUpdate: string;

	Name: string;
};
