package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social/ent"
	"social/internal/usecase"
)

type GroupHandler struct {
	groupUsecase *usecase.GroupUseCase
}

func NewGroupHandler() *GroupHandler {
	return &GroupHandler{}
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

	result, err := h.groupUsecase.GetGroupsRepo(c, u.ID)
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
