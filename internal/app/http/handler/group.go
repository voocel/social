package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social/ent"
	"social/internal/entity"
	"social/internal/usecase"
)

type GroupHandler struct {
	groupUsecase *usecase.GroupUseCase
}

func NewGroupHandler(g *usecase.GroupUseCase) *GroupHandler {
	return &GroupHandler{groupUsecase: g}
}

func (h *GroupHandler) GetGroups(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	u, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	result, err := h.groupUsecase.GetGroups(c, u.ID)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "ok"
	resp.Data = result
	c.JSON(http.StatusOK, resp)
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	u, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	req := &entity.Group{}
	if err := c.ShouldBind(req); err != nil {
		resp.Code = 1
		resp.Message = "params invalid"
		c.JSON(http.StatusOK, resp)
		return
	}
	req.Owner = u.ID

	result, err := h.groupUsecase.CreateGroup(c, req)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "ok"
	resp.Data = result
	c.JSON(http.StatusOK, resp)
}
