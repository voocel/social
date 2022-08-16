package http

import (
	"github.com/gin-gonic/gin"
	v1 "social/internal/delivery/http/v1"
)

func NewRouter(r *gin.Engine) {
	userRouter := r.Group("/v1/user")
	{
		userRouter.POST("/login", v1.Login)
		userRouter.POST("/register")
		userRouter.GET("/info")
	}

	friendRouter := r.Group("/v1/friend")
	{
		friendRouter.GET("/getFriends")
		friendRouter.POST("/addFriendApply")
		friendRouter.GET("/getFriendApply")
		friendRouter.GET("/agreeFriendApply")
		friendRouter.GET("/refuseFriendApply")
	}

	groupRouter := r.Group("/v1/group")
	{
		groupRouter.GET("/getGroups")
		groupRouter.POST("/createGroup")
	}
}
