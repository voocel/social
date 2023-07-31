package handler

import (
	"github.com/gin-gonic/gin"
	"social/internal/usecase"
)

type FriendApplyHandle struct {
	faUseCase *usecase.FriendApplyUseCase
}

func NewFriendApplyHandle(f *usecase.FriendApplyUseCase) *FriendApplyHandle {
	return &FriendApplyHandle{faUseCase: f}
}

func (h *FriendApplyHandle) AddFriendApply(c *gin.Context) {

}

func (h *FriendApplyHandle) GetFriendApply(c *gin.Context) {

}

func (h *FriendApplyHandle) AgreeFriendApply(c *gin.Context) {

}
