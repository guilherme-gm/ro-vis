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

export type MinI18n = {
	I18nId: number;
	LastUpdate: string;
	PtBrText: string;
};
