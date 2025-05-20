package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
)

type UpdatesController struct{}

func (ctlr *UpdatesController) List(c *gin.Context) {
	offset, err := queryInt(c, "start", 0)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// @TODO: Probably better make it a go routine
	patchRepo := repository.GetPatchRepository()
	count, err := patchRepo.GetUpdateCount(nil)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if count < int32(offset) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	updates, err := repository.GetPatchRepository().ListUpdates(nil, repository.Pagination{
		Offset: int32(offset),
		Limit:  100,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": updates})
}
