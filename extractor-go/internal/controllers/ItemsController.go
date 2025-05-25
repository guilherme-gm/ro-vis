package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

type ItemsController struct{}

type itemResponse struct {
	PreviousHistoryID       *int32
	HistoryID               *int32
	ItemID                  int32
	FileVersion             int32
	IdentifiedName          *string
	IdentifiedDescription   *string
	IdentifiedSprite        *string
	UnidentifiedName        *string
	UnidentifiedDescription *string
	UnidentifiedSprite      *string
	SlotCount               int8
	IsBook                  bool
	CanUseBuyingStore       bool
	CardPrefix              *string
	CardIsPostfix           bool
	CardIllustration        *string
	ClassNum                *int32
	IsCostume               bool
	EffectID                int32
	PackageID               int32
	MoveInfo                domain.ItemMoveInfo
}

type minItemResponse struct {
	ItemID         int32
	LastUpdate     string
	IdentifiedName *string
}

func (ctlr *ItemsController) formatFromTo(record repository.FromToRecord[domain.Item]) fromToRecordResponse[itemResponse] {
	var fromRec *recordResponse[itemResponse] = nil
	var toRec *recordResponse[itemResponse] = nil

	if record.From != nil {
		fromRec = &recordResponse[itemResponse]{
			Update: record.From.Update,
			Data:   nil,
		}
		fromRec.Data = &itemResponse{
			PreviousHistoryID:       sqlInt32ToPointer(record.From.Data.PreviousHistoryID),
			HistoryID:               sqlInt32ToPointer(record.From.Data.HistoryID),
			ItemID:                  record.From.Data.ItemID,
			FileVersion:             record.From.Data.FileVersion,
			IdentifiedName:          sqlStringToPointer(record.From.Data.IdentifiedName),
			IdentifiedDescription:   sqlStringToPointer(record.From.Data.IdentifiedDescription),
			IdentifiedSprite:        sqlStringToPointer(record.From.Data.IdentifiedSprite),
			UnidentifiedName:        sqlStringToPointer(record.From.Data.UnidentifiedName),
			UnidentifiedDescription: sqlStringToPointer(record.From.Data.UnidentifiedDescription),
			UnidentifiedSprite:      sqlStringToPointer(record.From.Data.UnidentifiedSprite),
			SlotCount:               record.From.Data.SlotCount,
			IsBook:                  record.From.Data.IsBook,
			CanUseBuyingStore:       record.From.Data.CanUseBuyingStore,
			CardPrefix:              sqlStringToPointer(record.From.Data.CardPrefix),
			CardIsPostfix:           record.From.Data.CardIsPostfix,
			CardIllustration:        sqlStringToPointer(record.From.Data.CardIllustration),
			ClassNum:                sqlInt32ToPointer(record.From.Data.ClassNum),
			IsCostume:               record.From.Data.IsCostume,
			EffectID:                record.From.Data.EffectID,
			PackageID:               record.From.Data.PackageID,
			MoveInfo:                record.From.Data.MoveInfo,
		}
	}

	if record.To != nil {
		toRec = &recordResponse[itemResponse]{
			Update: record.To.Update,
			Data:   nil,
		}
		toRec.Data = &itemResponse{
			PreviousHistoryID:       sqlInt32ToPointer(record.To.Data.PreviousHistoryID),
			HistoryID:               sqlInt32ToPointer(record.To.Data.HistoryID),
			ItemID:                  record.To.Data.ItemID,
			FileVersion:             record.To.Data.FileVersion,
			IdentifiedName:          sqlStringToPointer(record.To.Data.IdentifiedName),
			IdentifiedDescription:   sqlStringToPointer(record.To.Data.IdentifiedDescription),
			IdentifiedSprite:        sqlStringToPointer(record.To.Data.IdentifiedSprite),
			UnidentifiedName:        sqlStringToPointer(record.To.Data.UnidentifiedName),
			UnidentifiedDescription: sqlStringToPointer(record.To.Data.UnidentifiedDescription),
			UnidentifiedSprite:      sqlStringToPointer(record.To.Data.UnidentifiedSprite),
			SlotCount:               record.To.Data.SlotCount,
			IsBook:                  record.To.Data.IsBook,
			CanUseBuyingStore:       record.To.Data.CanUseBuyingStore,
			CardPrefix:              sqlStringToPointer(record.To.Data.CardPrefix),
			CardIsPostfix:           record.To.Data.CardIsPostfix,
			CardIllustration:        sqlStringToPointer(record.To.Data.CardIllustration),
			ClassNum:                sqlInt32ToPointer(record.To.Data.ClassNum),
			IsCostume:               record.To.Data.IsCostume,
			EffectID:                record.To.Data.EffectID,
			PackageID:               record.To.Data.PackageID,
			MoveInfo:                record.To.Data.MoveInfo,
		}
	}

	return fromToRecordResponse[itemResponse]{
		LastUpdated: sqlStringToPointer(record.LastUpdate),
		From:        fromRec,
		To:          toRec,
	}
}

func (ctlr *ItemsController) List(c *gin.Context) {
	offset, err := queryInt(c, "start", 0)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	itemRepo := repository.GetItemRepository()
	count, err := itemRepo.CountItems(nil)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if count < int32(offset) {
		fmt.Println("Out of range")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	items, err := itemRepo.GetItems(nil, repository.Pagination{
		Offset: int32(offset),
		Limit:  100,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	list := []minItemResponse{}
	for _, val := range items {
		list = append(list, minItemResponse{
			ItemID:         val.ItemID,
			LastUpdate:     val.Lastupdate,
			IdentifiedName: sqlStringToPointer(val.IdentifiedName),
		})
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": list})
}

func (ctlr *ItemsController) ListForUpdate(c *gin.Context) {
	offset, err := queryInt(c, "start", 0)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	update, err := paramAsUpdate(c, "update")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// @TODO: Probably better make it a go routine
	itemRepo := repository.GetItemRepository()
	count, err := itemRepo.CountChangesInUpdate(nil, update)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if count < offset {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	updates, err := itemRepo.GetChangesInUpdate(nil, update, repository.Pagination{
		Offset: int32(offset),
		Limit:  100,
	})
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	list := []fromToRecordResponse[itemResponse]{}
	for _, val := range updates {
		list = append(list, ctlr.formatFromTo(val))
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": list})
}

func (ctlr *ItemsController) ListForItem(c *gin.Context) {
	offset, err := queryInt(c, "start", 0)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	itemId, err := intParam(c, "itemId")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	itemRepo := repository.GetItemRepository()
	updates, err := itemRepo.GetItemHistory(nil, int32(itemId), repository.Pagination{
		Offset: int32(offset),
		Limit:  100,
	})
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	list := []fromToRecordResponse[itemResponse]{}
	for _, val := range updates {
		list = append(list, ctlr.formatFromTo(val))
	}

	c.JSON(http.StatusOK, gin.H{"total": len(list), "list": list})
}
