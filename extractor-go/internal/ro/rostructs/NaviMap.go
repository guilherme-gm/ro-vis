package rostructs

import (
	"strconv"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

// These covers LUA structures.
// TXT files are handled by loaders directly

type NaviMap struct {
	MapId       string `lua:"$$numeric:1"`
	Name        string `lua:"$$numeric:2"`
	SpecialCode int    `lua:"$$numeric:3"`
	LocationX   int    `lua:"$$numeric:4"`
	LocationY   int    `lua:"$$numeric:5"`
}

type NaviNpc struct {
	MapId      string `lua:"$$numeric:1"`
	UniqueCode int    `lua:"$$numeric:2"`
	Type       int    `lua:"$$numeric:3"`
	SpriteId   int    `lua:"$$numeric:4"`
	Name1      string `lua:"$$numeric:5"`
	Name2      string `lua:"$$numeric:6"`
	LocationX  int    `lua:"$$numeric:7"`
	LocationY  int    `lua:"$$numeric:8"`
}

func (n NaviNpc) ToDomain() domain.MapNpc {
	var npcType domain.NpcType
	switch n.Type {
	case 101:
		npcType = domain.NpcTypeDialog
	case 102:
		npcType = domain.NpcTypeShop
	default:
		panic("Invalid npc type: " + strconv.Itoa(n.Type))
	}

	return domain.MapNpc{
		Type:     npcType,
		SpriteId: n.SpriteId,
		Name1:    domain.NewLocalizableStringFromNavi(n.Name1),
		Name2:    n.Name2,
		Location: domain.MapCoord{
			MapId: n.MapId,
			X:     n.LocationX,
			Y:     n.LocationY,
		},
	}
}

type NaviMob struct {
	MapId      string `lua:"$$numeric:1"`
	UniqueCode int    `lua:"$$numeric:2"`
	Type       int    `lua:"$$numeric:3"`
	IdAmount   int    `lua:"$$numeric:4"`
	Name1      string `lua:"$$numeric:5"`
	Name2      string `lua:"$$numeric:6"`
	Level      int    `lua:"$$numeric:7"`
	Detail     int    `lua:"$$numeric:8"`
}

func (n NaviMob) ToDomain() domain.MapSpawn {
	spriteId := int(n.IdAmount & 0x0000FFFF)
	amount := int((int64(n.IdAmount) & 0xFFFF0000) >> 16)

	element := int((int64(n.Detail) & 0xFFFF0000) >> 16)
	size := int((n.Detail & 0x0000FF00) >> 8)
	race := int(n.Detail & 0x000000FF)

	var spawnType domain.SpawnType
	switch n.Type {
	case 300:
		spawnType = domain.SpawnTypeNormal
	case 301:
		spawnType = domain.SpawnTypeBoss
	default:
		panic("Invalid spawn type: " + strconv.Itoa(n.Type))
	}

	return domain.MapSpawn{
		Type:     spawnType,
		SpriteId: spriteId,
		Name1:    domain.NewLocalizableStringFromNavi(n.Name1),
		Name2:    n.Name2,
		Level:    n.Level,
		Amount:   amount,
		Element:  element,
		Size:     size,
		Race:     race,
	}
}

type NaviLink struct {
	MapId      string `lua:"$$numeric:1"`
	UniqueCode int    `lua:"$$numeric:2"`
	WarpType   int    `lua:"$$numeric:3"`
	SpriteId   int    `lua:"$$numeric:4"`
	Name1      string `lua:"$$numeric:5"`
	Name2      string `lua:"$$numeric:6"`
	FromX      int    `lua:"$$numeric:7"`
	FromY      int    `lua:"$$numeric:8"`
	ToMapId    string `lua:"$$numeric:9"`
	ToX        int    `lua:"$$numeric:10"`
	ToY        int    `lua:"$$numeric:11"`
}

func (n NaviLink) ToDomain() domain.MapWarp {
	var warpType domain.WarpType
	switch n.WarpType {
	case 200:
		warpType = domain.WarpTypeCommon
	case 201:
		warpType = domain.WarpTypeFreeNpc
	case 202:
		warpType = domain.WarpTypeKafraDts
	case 203:
		warpType = domain.WarpTypeCoolEventDts
	case 204:
		warpType = domain.WarpTypePaidNpc
	case 205:
		warpType = domain.WarpTypeAirport
	default:
		panic("Invalid warp type: " + strconv.Itoa(n.WarpType))
	}

	return domain.MapWarp{
		From: domain.MapCoord{
			MapId: n.MapId,
			X:     n.FromX,
			Y:     n.FromY,
		},
		To: domain.MapCoord{
			MapId: n.ToMapId,
			X:     n.ToX,
			Y:     n.ToY,
		},
		WarpType: warpType,
		SpriteId: n.SpriteId,
		Name1:    domain.NewLocalizableStringFromNavi(n.Name1),
		Name2:    n.Name2,
	}
}
