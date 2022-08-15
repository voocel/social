package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"social/internal/app/http/middleware"
	"social/internal/usecase"
	"social/internal/usecase/repo"
	"social/pkg/log"
	"social/pkg/postgres"
	"syscall"
	"time"
)

func Run() {
	r := gin.New()
	gin.SetMode(gin.DebugMode)
	pg, err := postgres.New("127.0.0.1:6008", postgres.MaxPoolSize(2))
	if err != nil {
		log.Fatal(fmt.Errorf("postgres New err: %w", err))
	}
	defer pg.Close()
	userUseCase := usecase.NewUserUseCase(repo.NewUserRepo(pg))
	r.Use(
		gin.Recovery(),
		middleware.JWTMiddleware(userUseCase, "prod", "social-key"),
	)
	NewRouter(r)

	srv := http.Server{
		Addr:    ":7788",
		Handler: r,
	}
	go func() {
		if err  = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		s := <-ch
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
		default:
			return
		}
	}
}
