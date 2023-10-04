package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social/ent"
	"social/internal/usecase"
)

type MessageHandler struct {
	mUseCase *usecase.MessageUseCase
}

func NewMessageHandler(m *usecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{mUseCase: m}
}

func (h *MessageHandler) GetMessage(c *gin.Context) {
	resp := new(ApiResponse)
	user, exists := c.Get("jwt-user")
	u, ok := user.(*ent.User)
	if !exists || !ok {
		resp.Code = 1
		resp.Message = "token invalid"
		c.JSON(http.StatusOK, resp)
		return
	}

	result, err := h.mUseCase.GetMessages(c, u.ID)
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
