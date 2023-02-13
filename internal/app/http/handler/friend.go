package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FriendHandler struct {
}

func NewFriendHandler() *FriendHandler {
	return &FriendHandler{}
}

func (f *FriendHandler) GetFriends(c *gin.Context) {
	resp := new(ApiResponse)
	uidStr, ok := c.GetQuery("uid")
	if !ok {
		resp.Code = 1
		resp.Message = "params invalid"
		c.JSON(http.StatusOK, resp)
		return
	}
	uid, _ := strconv.Atoi(uidStr)
	result := map[string]interface{}{
		"uid":       uid,
		"nickname":  "test",
		"sex":       1,
		"avatar":    "",
		"friend_id": 1,
	}
	resp.Message = "ok"
	resp.Data = []map[string]interface{}{result}
	c.JSON(http.StatusOK, resp)
}
