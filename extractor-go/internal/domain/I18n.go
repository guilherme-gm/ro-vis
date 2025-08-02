package domain

type I18n struct {
	PreviousHistoryID NullableInt64
	HistoryID         NullableInt64
	I18nId            uint64
	FileVersion       int32
	ContainerFile     string
	EnText            string
	PtBrText          string
	Active            bool
	Deleted           bool
}

func (i *I18n) Equals(otherI18n I18n) bool {
	return (i.I18nId == otherI18n.I18nId &&
		i.ContainerFile == otherI18n.ContainerFile &&
		i.EnText == otherI18n.EnText &&
		i.PtBrText == otherI18n.PtBrText &&
		i.Active == otherI18n.Active)
}
