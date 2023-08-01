package http

import (
	"github.com/gin-gonic/gin"
	"social/internal/app/http/handler"
	"social/internal/usecase"
)

type Router interface {
	Load(r *gin.Engine)
}

func getRouters(u *usecase.UserUseCase, f *usecase.FriendUseCase, fa *usecase.FriendApplyUseCase, g *usecase.GroupUseCase) (routers []Router) {
	userHandler := handler.NewUserHandler(u)
	ur := newUserRouter(userHandler)

	friendHandler := handler.NewFriendHandler(f)
	fr := newFriendRouter(friendHandler)

	friendApplyHandler := handler.NewFriendApplyHandle(fa)
	far := newFriendApplyRouter(friendApplyHandler)

	groupHandler := handler.NewGroupHandler(g)
	gr := newGroupRouter(groupHandler)

	routers = append(routers, ur, fr, far, gr)
	return
}

// ######################### User Router #########################
type userRouter struct {
	h *handler.UserHandler
}

func newUserRouter(h *handler.UserHandler) *userRouter {
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
	h *handler.FriendHandler
}

func newFriendRouter(h *handler.FriendHandler) *friendRouter {
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

// ######################### Friend Router #########################
type friendApplyRouter struct {
	h *handler.FriendApplyHandler
}

func newFriendApplyRouter(h *handler.FriendApplyHandler) *friendApplyRouter {
	return &friendApplyRouter{h: h}
}

func (r *friendApplyRouter) Load(g *gin.Engine) {
	fr := g.Group("/v1/friend")
	{
		fr.GET("/getFriends", r.h.GetFriendApply)
		fr.POST("/addFriendApply")
		fr.GET("/getFriendApply")
		fr.GET("/agreeFriendApply")
		fr.GET("/refuseFriendApply")
	}
}

// ######################### Group Router #########################
type groupRouter struct {
	h *handler.GroupHandler
}

func newGroupRouter(h *handler.GroupHandler) *groupRouter {
	return &groupRouter{h: h}
}

func (r *groupRouter) Load(g *gin.Engine) {
	gr := g.Group("/v1/group")
	{
		gr.GET("/getGroups", r.h.GetGroups)
		gr.POST("/createGroup")
	}
}
