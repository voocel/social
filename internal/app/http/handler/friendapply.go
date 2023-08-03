package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social/ent"
	"social/internal/entity"
	"social/internal/usecase"
	"strconv"
)

type FriendApplyHandler struct {
	faUseCase *usecase.FriendApplyUseCase
}

func NewFriendApplyHandle(f *usecase.FriendApplyUseCase) *FriendApplyHandler {
	return &FriendApplyHandler{faUseCase: f}
}

func (h *FriendApplyHandler) AddFriendApply(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	u, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	req := &entity.FriendApply{}
	if err := c.ShouldBind(req); err != nil {
		resp.Code = 1
		resp.Message = "params invalid"
		c.JSON(http.StatusOK, resp)
		return
	}
	req.FromId = u.ID

	if _, err := h.faUseCase.AddFriendApply(c, req); err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "ok"
	c.JSON(http.StatusOK, resp)
}

func (h *FriendApplyHandler) GetFriendApply(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	u, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	result, err := h.faUseCase.GetFriendApply(c, u.ID)
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

func (h *FriendApplyHandler) AgreeFriendApply(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	u, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}
	applyId, ok := c.GetQuery("apply_id")
	if !ok {
		resp.Code = 1
		resp.Message = "apply id not exist"
		c.JSON(http.StatusOK, resp)
		return
	}
	id, _ := strconv.Atoi(applyId)

	if _, err := h.faUseCase.AgreeFriendApply(c, int64(id), u.ID); err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Message = "ok"
	c.JSON(http.StatusOK, resp)
}
