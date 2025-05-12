package domain

import (
	"database/sql"
)

type ItemMoveInfo struct {
	CanDrop               bool
	CanTrade              bool
	CanMoveToStorage      bool
	CanMoveToCart         bool
	CanSellToNpc          bool
	CanMail               bool
	CanAuction            bool
	CanMoveToGuildStorage bool
	CommentName           string
}

func NewItemMoveInfo() ItemMoveInfo {
	return ItemMoveInfo{
		CanDrop:               true,
		CanTrade:              true,
		CanMoveToStorage:      true,
		CanMoveToCart:         true,
		CanSellToNpc:          true,
		CanMail:               true,
		CanAuction:            true,
		CanMoveToGuildStorage: true,
	}
}

func (i *ItemMoveInfo) Equals(other ItemMoveInfo) bool {
	return (i.CanDrop == other.CanDrop &&
		i.CanTrade == other.CanTrade &&
		i.CanMoveToStorage == other.CanMoveToStorage &&
		i.CanMoveToCart == other.CanMoveToCart &&
		i.CanSellToNpc == other.CanSellToNpc &&
		i.CanMail == other.CanMail &&
		i.CanAuction == other.CanAuction &&
		i.CanMoveToGuildStorage == other.CanMoveToGuildStorage &&
		i.CommentName == other.CommentName)
}

type Item struct {
	PreviousHistoryID       sql.NullInt32
	HistoryID               sql.NullInt32
	ItemID                  int32
	FileVersion             int32
	IdentifiedName          sql.NullString
	IdentifiedDescription   sql.NullString
	IdentifiedSprite        sql.NullString
	UnidentifiedName        sql.NullString
	UnidentifiedDescription sql.NullString
	UnidentifiedSprite      sql.NullString
	SlotCount               int8
	IsBook                  bool
	CanUseBuyingStore       bool
	CardPrefix              sql.NullString
	CardIsPostfix           bool
	CardIllustration        sql.NullString
	ClassNum                sql.NullInt32
	MoveInfo                ItemMoveInfo
	Deleted                 bool
}

func NewItem(itemID int32, fileVersion int32) Item {
	return Item{
		ItemID:      itemID,
		FileVersion: fileVersion,
		MoveInfo:    NewItemMoveInfo(),
	}
}

func (i *Item) Equals(otherItem Item) bool {
	// FileVersion is not checked, if the file has changed but the Item is the same, we don't care.
	return (i.ItemID == otherItem.ItemID &&
		i.IdentifiedName == otherItem.IdentifiedName &&
		i.IdentifiedDescription == otherItem.IdentifiedDescription &&
		i.IdentifiedSprite == otherItem.IdentifiedSprite &&
		i.UnidentifiedName == otherItem.UnidentifiedName &&
		i.UnidentifiedDescription == otherItem.UnidentifiedDescription &&
		i.UnidentifiedSprite == otherItem.UnidentifiedSprite &&
		i.SlotCount == otherItem.SlotCount &&
		i.IsBook == otherItem.IsBook &&
		i.CanUseBuyingStore == otherItem.CanUseBuyingStore &&
		i.CardPrefix == otherItem.CardPrefix &&
		i.CardIsPostfix == otherItem.CardIsPostfix &&
		i.CardIllustration == otherItem.CardIllustration &&
		i.ClassNum == otherItem.ClassNum &&
		i.MoveInfo.Equals(otherItem.MoveInfo))
}
