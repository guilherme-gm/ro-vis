package domain

type MapCoord struct {
	MapId string
	X     int
	Y     int
}

type Map struct {
	PreviousHistoryID NullableInt32
	HistoryID         NullableInt32
	Id                string
	FileVersion       int32
	Name              NullableString
	SpecialCode       NullableInt32
	Mp3Name           NullableString
	Npcs              []MapNpc
	Warps             []MapWarp
	Spawns            []MapSpawn
}

type NpcType string

const (
	NpcTypeDialog NpcType = "dialog"
	NpcTypeShop   NpcType = "shop"
)

type MapNpc struct {
	Type     NpcType
	SpriteId int
	Name1    LocalizableString
	Name2    string
	Location MapCoord
}

type SpawnType string

const (
	SpawnTypeNormal SpawnType = "normal"
	SpawnTypeBoss   SpawnType = "boss"
)

type MapSpawn struct {
	Type     SpawnType
	SpriteId int
	Name1    LocalizableString
	Name2    string
	Level    int
	Amount   int
	Element  int
	Size     int
	Race     int
}

type MapWarp struct {
	From     MapCoord
	To       MapCoord
	WarpType WarpType
	SpriteId int
	Name1    LocalizableString
	Name2    string
}

type WarpType string

const (
	WarpTypeCommon       WarpType = "common"
	WarpTypeFreeNpc      WarpType = "free_npc"
	WarpTypeKafraDts     WarpType = "kafra_dts"
	WarpTypeCoolEventDts WarpType = "cool_event_dts"
	WarpTypePaidNpc      WarpType = "paid_npc"
	WarpTypeAirport      WarpType = "airport"
)
