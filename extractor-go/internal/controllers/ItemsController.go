package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
)

type ItemsController struct{}

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

	c.JSON(http.StatusOK, gin.H{"total": count, "list": items})
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

	c.JSON(http.StatusOK, gin.H{"total": count, "list": updates})
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

	c.JSON(http.StatusOK, gin.H{"total": len(updates), "list": updates})
}
