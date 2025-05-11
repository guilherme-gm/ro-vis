package rostructs

// Note: Item V1 does not exists because it was many TXT files, so ItemV1Parser handled it

/**
 * LUA part of ItemV2 (TXT part is handled by ItemV2Parser directly)
 */
type ItemV2 struct {
	ItemID                      int      `lua:"@index"`
	UnidentifiedDisplayName     string   `lua:"unidentifiedDisplayName"`
	UnidentifiedResourceName    string   `lua:"unidentifiedResourceName"`
	UnidentifiedDescriptionName []string `lua:"unidentifiedDescriptionName"`
	IdentifiedDisplayName       string   `lua:"identifiedDisplayName"`
	IdentifiedResourceName      string   `lua:"identifiedResourceName"`
	IdentifiedDescriptionName   []string `lua:"identifiedDescriptionName"`
	SlotCount                   int      `lua:"slotCount"`
	ClassNum                    int
}
