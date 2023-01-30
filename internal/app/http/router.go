package http

import (
	"github.com/gin-gonic/gin"
	v1 "social/internal/delivery/http/v1"
	"social/internal/usecase"
)

type Router interface {
	Load(r *gin.Engine)
}

func getRouters(userUsecase *usecase.UserUseCase) (routers []Router) {
	userHandler := v1.NewUserHandler(userUsecase)
	ur := NewUserRouter(userHandler)

	friendHandler := v1.NewFriendHandler()
	fr := NewFriendRouter(friendHandler)

	groupHandler := v1.NewGroupHandler()
	gr := NewGroupRouter(groupHandler)

	routers = append(routers, ur, fr, gr)
	return
}

// ######################### User Router #########################
type userRouter struct {
	h *v1.UserHandler
}

func NewUserRouter(h *v1.UserHandler) *userRouter {
	return &userRouter{h: h}
}

func (r *userRouter) Load(g *gin.Engine) {
	ur := g.Group("/v1/user")
	{
		ur.POST("/login", r.h.Login)
		ur.POST("/register", r.h.Register)
		ur.GET("/info", r.h.Info)
	}
}

// ######################### Friend Router #########################
type friendRouter struct {
	h *v1.FriendHandler
}

func NewFriendRouter(h *v1.FriendHandler) *friendRouter {
	return &friendRouter{h: h}
}

func (r *friendRouter) Load(g *gin.Engine) {
	fr := g.Group("/v1/friend")
	{
		fr.GET("/getFriends", r.h.GetFriends)
		fr.POST("/addFriendApply")
		fr.GET("/getFriendApply")
		fr.GET("/agreeFriendApply")
		fr.GET("/refuseFriendApply")
	}
}

// ######################### Group Router #########################
type groupRouter struct {
	h *v1.GroupHandler
}

func NewGroupRouter(h *v1.GroupHandler) *groupRouter {
	return &groupRouter{h: h}
}

func (r *groupRouter) Load(g *gin.Engine) {
	gr := g.Group("/v1/group")
	{
		gr.GET("/getGroups")
		gr.POST("/createGroup")
	}
}
