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
	fUseCase  *usecase.FriendUseCase
}

func NewFriendApplyHandle(f *usecase.FriendUseCase, fa *usecase.FriendApplyUseCase) *FriendApplyHandler {
	return &FriendApplyHandler{faUseCase: fa, fUseCase: f}
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

	req := &entity.FriendApplyReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Code = 1
		resp.Message = "params invalid"
		c.JSON(http.StatusOK, resp)
		return
	}
	args := &entity.FriendApply{}
	args.FromId = u.ID
	args.ToId = req.FriendId
	args.Remark = req.ApplyInfo

	if _, err := h.faUseCase.AddFriendApply(c, args); err != nil {
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
	fromId, ok := c.GetQuery("from_id")
	if !ok {
		resp.Code = 1
		resp.Message = "apply id not exist"
		c.JSON(http.StatusOK, resp)
		return
	}
	fid, _ := strconv.Atoi(fromId)
	tid := u.ID

	if _, err := h.faUseCase.AgreeFriendApply(c, int64(fid), tid); err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	req := &entity.Friend{
		Uid:      tid,
		FriendId: int64(fid),
	}
	if _, err := h.fUseCase.AddFriend(c, req); err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "ok"
	c.JSON(http.StatusOK, resp)
}
