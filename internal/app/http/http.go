package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"social/ent"
	"social/internal/app/http/handler"
	"social/internal/app/http/middleware"
	"social/internal/usecase"
	"social/internal/usecase/repo"
	"social/pkg/log"

	_ "github.com/lib/pq"
	_ "social/docs"
)

type Server struct {
	ent *ent.Client
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

	s.initEnt()
	userUseCase := usecase.NewUserUseCase(repo.NewUserRepo(s.ent))
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
	s.routerLoad(g, getRouters(userUseCase)...)

	srv := http.Server{
		Addr:    viper.GetString("http.addr"),
		Handler: g,
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		s := <-ch
		log.Infof("[%v]Shutting down...", s)
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			if err = srv.Shutdown(ctx); err != nil {
				panic(err)
			}
			log.Sync()
			cancel()
			return
		case syscall.SIGHUP:
			return
		default:
			return
		}
	}
}

func (s *Server) initEnt() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.database"),
		viper.GetString("postgres.sslmode"),
	)

	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatal(err)
	}

	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
			meta, ok := ctx.Value(middleware.AuditID).(middleware.Event)
			if !ok {
				return next.Mutate(ctx, mutation)
			}

			val, err := next.Mutate(ctx, mutation)

			meta.Table = mutation.Type()
			meta.Action = middleware.Action(mutation.Op().String())

			newValues, _ := json.Marshal(val)
			meta.NewValues = string(newValues)
			log.Info(meta)

			return val, err
		})
	})

	s.ent = client
}
