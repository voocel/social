package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social/ent"
	"social/internal/usecase"
)

type FriendHandler struct {
	friendUsecase *usecase.FriendUseCase
}

func NewFriendHandler(f *usecase.FriendUseCase) *FriendHandler {
	return &FriendHandler{friendUsecase: f}
}

func (f *FriendHandler) GetFriends(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	u, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	result, err := f.friendUsecase.GetFriendsRepo(c, u.ID)
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
