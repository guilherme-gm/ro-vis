package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

type MapsController struct{}

type ListMapsParams struct {
	Query PaginateQuery
}

func (ctlr *MapsController) List(c *gin.Context, params ListMapsParams) {
	mapRepo := c.MustGet("x-server").(*server.Server).Repositories.MapRepository
	count, err := mapRepo.CountMaps(nil)
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch count", err))
		return
	}

	if count < int32(params.Query.Start) {
		c.Error(NewBadRequestError("offset is out of range", nil))
		return
	}

	maps, err := mapRepo.GetMaps(nil, repository.Pagination{
		Offset: int32(params.Query.Start),
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch maps", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": maps})
}

type ListMapsForUpdateParams struct {
	Params struct {
		Update string `uri:"update" binding:"updateStr"`
	}
	Query PaginateQuery
}

func (ctlr *MapsController) ListForUpdate(c *gin.Context, params ListMapsForUpdateParams) {
	// @TODO: Probably better make it a go routine
	mapRepo := c.MustGet("x-server").(*server.Server).Repositories.MapRepository
	count, err := mapRepo.CountChangesInUpdate(nil, params.Params.Update)
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch count", err))
		return
	}

	if count < int(params.Query.Start) {
		c.Error(NewBadRequestError("offset is out of range", nil))
		return
	}

	updates, err := mapRepo.GetChangesInUpdate(nil, params.Params.Update, repository.Pagination{
		Offset: params.Query.Start,
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch maps", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": updates})
}

type ListForMapParams struct {
	Params struct {
		MapId string `uri:"mapId" binding:"min=1"`
	}
	Query PaginateQuery
}

func (ctlr *MapsController) ListForItem(c *gin.Context, params ListForMapParams) {
	mapRepo := c.MustGet("x-server").(*server.Server).Repositories.MapRepository
	updates, err := mapRepo.GetMapHistory(nil, params.Params.MapId, repository.Pagination{
		Offset: params.Query.Start,
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch maps", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": len(updates), "list": updates})
}
