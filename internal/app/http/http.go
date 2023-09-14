package http

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"social/config"
	"social/internal/app/http/handler"
	"social/internal/app/http/middleware"
	"social/internal/usecase"
	"social/internal/usecase/repo"
	"social/pkg/database/ent"

	_ "social/docs"
)

type Server struct {
	srv http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) routerLoad(g *gin.Engine, rs ...Router) {
	for _, r := range rs {
		r.Load(g)
	}
}

func (s *Server) Run() {
	var err error
	g := gin.New()
	gin.SetMode(gin.DebugMode)

	cfg := ent.PgConfig{}
	copier.Copy(&cfg, config.Conf.Postgres)
	entClient := ent.InitEnt(cfg)

	userUseCase := usecase.NewUserUseCase(repo.NewUserRepo(entClient))
	friendUseCase := usecase.NewFriendUseCase(repo.NewFriendRepo(entClient))
	friendApplyUseCase := usecase.NewFriendApplyUseCase(repo.NewFriendApplyRepo(entClient))
	groupUseCase := usecase.NewGroupUseCase(repo.NewGroupRepo(entClient))

	g.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.Logger,
		middleware.CorsMiddleware(),
		middleware.JWTMiddleware(userUseCase),
	)
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "404 not found!")
	})
	g.GET("/ping", handler.Ping())
	g.Group("/swagger").GET("*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.routerLoad(g, getRouters(userUseCase, friendUseCase, friendApplyUseCase, groupUseCase)...)

	srv := http.Server{
		Addr:    viper.GetString("http.addr"),
		Handler: g,
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
