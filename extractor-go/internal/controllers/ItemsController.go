package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
)

type ItemsController struct{}

type ListItemsParams struct {
	Query PaginateQuery
}

func (ctlr *ItemsController) List(c *gin.Context, params ListItemsParams) {
	itemRepo := repository.GetItemRepository()
	count, err := itemRepo.CountItems(nil)
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch count", err))
		return
	}

	if count < int32(params.Query.Start) {
		c.Error(NewBadRequestError("offset is out of range", nil))
		return
	}

	items, err := itemRepo.GetItems(nil, repository.Pagination{
		Offset: int32(params.Query.Start),
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch items", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": items})
}

type ListItemsForUpdateParams struct {
	Params struct {
		Update string `uri:"update" binding:"updateStr"`
	}
	Query PaginateQuery
}

func (ctlr *ItemsController) ListForUpdate(c *gin.Context, params ListItemsForUpdateParams) {
	// @TODO: Probably better make it a go routine
	itemRepo := repository.GetItemRepository()
	count, err := itemRepo.CountChangesInUpdate(nil, params.Params.Update)
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch count", err))
		return
	}

	if count < int(params.Query.Start) {
		c.Error(NewBadRequestError("offset is out of range", nil))
		return
	}

	updates, err := itemRepo.GetChangesInUpdate(nil, params.Params.Update, repository.Pagination{
		Offset: params.Query.Start,
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch items", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": updates})
}

type ListForItemParams struct {
	Params struct {
		ItemId int32 `uri:"itemId" binding:"min=1"`
	}
	Query PaginateQuery
}

func (ctlr *ItemsController) ListForItem(c *gin.Context, params ListForItemParams) {
	itemRepo := repository.GetItemRepository()
	updates, err := itemRepo.GetItemHistory(nil, params.Params.ItemId, repository.Pagination{
		Offset: params.Query.Start,
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch items", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": len(updates), "list": updates})
}
