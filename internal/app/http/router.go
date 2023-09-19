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

	friendHandler := handler.NewFriendHandler(f, fa)
	fr := newFriendRouter(friendHandler)

	groupHandler := handler.NewGroupHandler(g)
	gr := newGroupRouter(groupHandler)

	routers = append(routers, ur, fr, gr)
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
		ur.GET("/getEmoji", r.h.GetEmoji)
		ur.GET("/updateAvatar", r.h.UpdateAvatar)
		ur.POST("/uploadFile", r.h.UploadFile)
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
		fr.GET("/list", r.h.GetFriends)
		fr.POST("/addApply", r.h.AddFriendApply)
		fr.GET("/getApply", r.h.GetFriendApply)
		fr.GET("/agreeApply", r.h.AgreeFriendApply)
		fr.GET("/refuseApply", r.h.RefuseFriendApply)
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
		gr.GET("/list", r.h.GetGroups)
		gr.POST("/create")
	}
}
