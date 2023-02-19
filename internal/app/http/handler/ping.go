package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "\r\nSUCCESS")
	}
}
