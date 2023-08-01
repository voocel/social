package handler

import (
	"github.com/gin-gonic/gin"
	"social/internal/usecase"
)

type FriendApplyHandler struct {
	faUseCase *usecase.FriendApplyUseCase
}

func NewFriendApplyHandle(f *usecase.FriendApplyUseCase) *FriendApplyHandler {
	return &FriendApplyHandler{faUseCase: f}
}

func (h *FriendApplyHandler) AddFriendApply(c *gin.Context) {

}

func (h *FriendApplyHandler) GetFriendApply(c *gin.Context) {

}

func (h *FriendApplyHandler) AgreeFriendApply(c *gin.Context) {

}
