export type I18n = {
	PreviousHistoryID: number | null;
	HistoryID: number | null;
	I18nId: number;
	FileVersion: number;
	ContainerFile: string;
	EnText: string;
	PtBrText: string;
	Active: boolean;
	Deleted: boolean;
};

export type I18nRecord = {
	update: string;
	data: I18n;
};

export type I18nFromToRecord = {
	from?: I18nRecord;
	to: I18nRecord;
	lastUpdate?: string;
};

export type MinI18n = {
	I18nId: number;
	LastUpdate: string;
	PtBrText: string;
};
