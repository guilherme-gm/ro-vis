package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

type I18nController struct{}

type ListI18nParams struct {
	Query PaginateQuery
}

func (ctlr *I18nController) List(c *gin.Context, params ListI18nParams) {
	i18nRepo := c.MustGet("x-server").(*server.Server).Repositories.I18nRepository
	count, err := i18nRepo.CountI18ns(nil)
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch count", err))
		return
	}

	if int32(count) < params.Query.Start {
		c.Error(NewBadRequestError("offset is out of range", nil))
		return
	}

	i18ns, err := i18nRepo.GetI18ns(nil, repository.Pagination{
		Offset: params.Query.Start,
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch i18n records", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": i18ns})
}

type ListI18nForUpdateParams struct {
	Params struct {
		Update string `uri:"update" binding:"updateStr"`
	}
	Query PaginateQuery
}

func (ctlr *I18nController) ListForUpdate(c *gin.Context, params ListI18nForUpdateParams) {
	i18nRepo := c.MustGet("x-server").(*server.Server).Repositories.I18nRepository

	// Get count of changes in this update
	count, err := i18nRepo.CountChangesInUpdate(nil, params.Params.Update)
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch count", err))
		return
	}

	if int32(count) < params.Query.Start {
		c.Error(NewBadRequestError("offset is out of range", nil))
		return
	}

	// Get the changes
	changes, err := i18nRepo.GetChangesInUpdate(nil, params.Params.Update, repository.Pagination{
		Offset: params.Query.Start,
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch i18n changes", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": changes})
}

type ListForI18nParams struct {
	Params struct {
		I18nId string `uri:"i18nId" binding:"min=1"`
	}
	Query PaginateQuery
}

func (ctlr *I18nController) ListForI18n(c *gin.Context, params ListForI18nParams) {
	i18nRepo := c.MustGet("x-server").(*server.Server).Repositories.I18nRepository

	// Get the history for this i18n record
	history, err := i18nRepo.GetI18nHistory(nil, params.Params.I18nId, repository.Pagination{
		Offset: params.Query.Start,
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch i18n history", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": len(history), "list": history})
}
