package http

import (
	"github.com/gin-gonic/gin"
	"social/ent"
	"social/internal/app/http/handler"
	"social/internal/usecase"
	"social/internal/usecase/repo"
)

type Router interface {
	Load(r *gin.Engine)
}

func getRouters(entClient *ent.Client) (routers []Router) {
	u := usecase.NewUserUseCase(repo.NewUserRepo(entClient))
	f := usecase.NewFriendUseCase(repo.NewFriendRepo(entClient))
	fa := usecase.NewFriendApplyUseCase(repo.NewFriendApplyRepo(entClient))
	g := usecase.NewGroupUseCase(repo.NewGroupRepo(entClient))
	gm := usecase.NewGroupMemberUseCase(repo.NewGroupMemberRepo(entClient))
	m := usecase.NewMessageUseCase(repo.NewMessageRepo(entClient))

	userHandler := handler.NewUserHandler(u)
	ur := newUserRouter(userHandler)

	friendHandler := handler.NewFriendHandler(f, fa)
	fr := newFriendRouter(friendHandler)

	groupHandler := handler.NewGroupHandler(g, gm)
	gr := newGroupRouter(groupHandler)

	messageHandler := handler.NewMessageHandler(m)
	mr := newMessageRouter(messageHandler)

	routers = append(routers, ur, fr, gr, mr)
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
		ur.PUT("/updateAvatar", r.h.UpdateAvatar)
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
		fr.PUT("/agreeApply", r.h.AgreeFriendApply)
		fr.PUT("/refuseApply", r.h.RefuseFriendApply)
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
		gr.POST("/create", r.h.CreateGroup)
		gr.POST("/joinGroup", r.h.JoinGroup)
	}
}

// ######################### Message Router #########################
type messageRouter struct {
	h *handler.MessageHandler
}

func newMessageRouter(h *handler.MessageHandler) *messageRouter {
	return &messageRouter{h: h}
}

func (r *messageRouter) Load(g *gin.Engine) {
	mr := g.Group("/v1/message")
	{
		mr.GET("/list", r.h.GetMessage)
	}
}
