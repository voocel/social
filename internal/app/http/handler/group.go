package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"social/ent"
	"social/internal/entity"
	"social/internal/usecase"
)

type GroupHandler struct {
	groupUsecase       *usecase.GroupUseCase
	groupMemberUsecase *usecase.GroupMemberUseCase
}

func NewGroupHandler(g *usecase.GroupUseCase, gm *usecase.GroupMemberUseCase) *GroupHandler {
	return &GroupHandler{groupUsecase: g, groupMemberUsecase: gm}
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
	req.CreatedUid = u.ID

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

	result, err := h.groupMemberUsecase.GetGroups(c, u.ID)
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

func (h *GroupHandler) JoinGroup(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	u, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	req := &entity.JoinGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Code = 1
		resp.Message = "params invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	_, err := h.groupUsecase.GetGroupById(c, req.GroupId)
	if err != nil || ent.IsNotFound(err) {
		resp.Code = 1
		resp.Message = fmt.Sprintf("Group not exists: %v", req.GroupId)
		c.JSON(http.StatusOK, resp)
		return
	}

	b, err := h.groupMemberUsecase.ExistsGroupMember(c, u.ID, req.GroupId)
	if err != nil || b {
		resp.Code = 1
		resp.Message = "You are already in this group"
		c.JSON(http.StatusOK, resp)
		return
	}

	group := &ent.GroupMember{
		UID:     u.ID,
		GroupID: req.GroupId,
		Remark:  req.Remark,
	}
	_, err = h.groupMemberUsecase.CreateGroupMember(c, group)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Message = "ok"
	c.JSON(http.StatusOK, resp)
}
