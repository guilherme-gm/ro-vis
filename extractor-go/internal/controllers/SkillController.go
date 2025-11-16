package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

type SkillController struct{}

type ListSkillParams struct {
	Query PaginateQuery
}

func (ctlr *SkillController) List(c *gin.Context, params ListSkillParams) {
	skillRepo := c.MustGet("x-server").(*server.Server).Repositories.SkillRepository
	count, err := skillRepo.CountSkills(nil)
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch count", err))
		return
	}

	if count < int32(params.Query.Start) {
		c.Error(NewBadRequestError("offset is out of range", nil))
		return
	}

	Skills, err := skillRepo.GetSkills(nil, repository.Pagination{
		Offset: int32(params.Query.Start),
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch Skills", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": count, "list": Skills})
}

type ListForSkillParams struct {
	Params struct {
		SkillId int32 `uri:"skillId" binding:"min=1"`
	}
	Query PaginateQuery
}

func (ctlr *SkillController) ListForItem(c *gin.Context, params ListForSkillParams) {
	itemRepo := c.MustGet("x-server").(*server.Server).Repositories.SkillRepository
	updates, err := itemRepo.GetSkillHistory(nil, params.Params.SkillId, repository.Pagination{
		Offset: params.Query.Start,
		Limit:  100,
	})
	if err != nil {
		c.Error(NewInternalServerError("failed to fetch Skill", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": len(updates), "list": updates})
}
