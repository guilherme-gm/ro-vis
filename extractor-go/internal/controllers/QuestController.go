package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

type QuestController struct{}

type ListQuestParams struct {
	Query PaginateQuery
}

func (ctlr *QuestController) List(c *gin.Context, params ListQuestParams) {
	questRepo := c.MustGet("x-server").(*server.Server).Repositories.QuestRepository
	count, err := questRepo.CountQuests(nil)
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch count", err))
		return
	}

	if count < int32(params.Query.Start) {
		c.Error(NewBadRequestError("offset is out of range", nil))
		return
	}

	Quest, err := questRepo.GetQuests(nil, repository.Pagination{
		Offset: int32(params.Query.Start),
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch Quest", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": Quest})
}

type ListQuestForUpdateParams struct {
	Params struct {
		Update string `uri:"update" binding:"updateStr"`
	}
	Query PaginateQuery
}

func (ctlr *QuestController) ListForUpdate(c *gin.Context, params ListQuestForUpdateParams) {
	// @TODO: Probably better make it a go routine
	itemRepo := c.MustGet("x-server").(*server.Server).Repositories.QuestRepository
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
		c.Error(NewInternalServerError("failed to fetch Quest", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": updates})
}

type ListForQuestParams struct {
	Params struct {
		QuestId int32 `uri:"questId" binding:"min=1"`
	}
	Query PaginateQuery
}

func (ctlr *QuestController) ListForItem(c *gin.Context, params ListForQuestParams) {
	itemRepo := c.MustGet("x-server").(*server.Server).Repositories.QuestRepository
	updates, err := itemRepo.GetQuestHistory(nil, params.Params.QuestId, repository.Pagination{
		Offset: params.Query.Start,
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch Quest", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": len(updates), "list": updates})
}
