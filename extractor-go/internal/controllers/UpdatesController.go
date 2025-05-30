package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
)

type UpdatesController struct{}

type ListUpdatesParams struct {
	Query PaginateQuery
}

func (ctlr *UpdatesController) List(c *gin.Context, params ListUpdatesParams) {
	// @TODO: Probably better make it a go routine
	patchRepo := repository.GetPatchRepository()
	count, err := patchRepo.GetUpdateCount(nil)
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch count", err))
		return
	}

	if count < int32(params.Query.Start) {
		c.Error(NewBadRequestError("offset is out of range", nil))
		return
	}

	updates, err := repository.GetPatchRepository().ListUpdates(nil, repository.Pagination{
		Offset: int32(params.Query.Start),
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch updates", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": updates})
}
