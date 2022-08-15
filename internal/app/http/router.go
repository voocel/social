package http

import (
	"github.com/gin-gonic/gin"
	v1 "social/internal/delivery/http/v1"
)

func NewRouter(r *gin.Engine)  {
	group := r.Group("/v1")
	{
		group.GET("/")
		group.POST("/login", v1.Login)
	}
}