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

/**
 * LUA part of ItemV5 (TXT part is handled by ItemV5Parser directly)
 */
type ItemV5 struct {
	ItemID                      int      `lua:"@index"`
	UnidentifiedDisplayName     string   `lua:"unidentifiedDisplayName"`
	UnidentifiedResourceName    string   `lua:"unidentifiedResourceName"`
	UnidentifiedDescriptionName []string `lua:"unidentifiedDescriptionName"`
	IdentifiedDisplayName       string   `lua:"identifiedDisplayName"`
	IdentifiedResourceName      string   `lua:"identifiedResourceName"`
	IdentifiedDescriptionName   []string `lua:"identifiedDescriptionName"`
	SlotCount                   int      `lua:"slotCount"`
	ClassNum                    int
	Costume                     bool `lua:"costume"`
}

/**
 * LUA part of ItemV6 (TXT part is handled by ItemV6Parser directly)
 */
type ItemV6 struct {
	ItemID                      int      `lua:"@index"`
	UnidentifiedDisplayName     string   `lua:"unidentifiedDisplayName"`
	UnidentifiedResourceName    string   `lua:"unidentifiedResourceName"`
	UnidentifiedDescriptionName []string `lua:"unidentifiedDescriptionName"`
	IdentifiedDisplayName       string   `lua:"identifiedDisplayName"`
	IdentifiedResourceName      string   `lua:"identifiedResourceName"`
	IdentifiedDescriptionName   []string `lua:"identifiedDescriptionName"`
	SlotCount                   int      `lua:"slotCount"`
	ClassNum                    int
	Costume                     bool `lua:"costume"`
	EffectID                    int
}

/**
 * LUA part of ItemV7 (TXT part is handled by ItemV7Parser directly)
 */
type ItemV7 struct {
	ItemID                      int      `lua:"@index"`
	UnidentifiedDisplayName     string   `lua:"unidentifiedDisplayName"`
	UnidentifiedResourceName    string   `lua:"unidentifiedResourceName"`
	UnidentifiedDescriptionName []string `lua:"unidentifiedDescriptionName"`
	IdentifiedDisplayName       string   `lua:"identifiedDisplayName"`
	IdentifiedResourceName      string   `lua:"identifiedResourceName"`
	IdentifiedDescriptionName   []string `lua:"identifiedDescriptionName"`
	SlotCount                   int      `lua:"slotCount"`
	ClassNum                    int
	Costume                     bool `lua:"costume"`
	EffectID                    int
	PackageID                   int
}
