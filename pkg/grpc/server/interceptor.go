package server

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerPrometheus() grpc.UnaryServerInterceptor {
	return grpc_prometheus.UnaryServerInterceptor
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	head := metadata.Pairs("action", "123")
	_ = grpc.SetHeader(ctx, head)

	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic err: %v; stack: %v", err, string(debug.Stack()))
		} else {
			log.Printf("grpc invoker method: %+v, request:%+v, response:%+v, md: %v, duration: %v", info.FullMethod, req, resp, md, time.Since(start))
		}
	}()
	return handler(ctx, req)
}
