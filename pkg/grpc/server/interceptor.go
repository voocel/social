package server

import (
	"context"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerPrometheus() grpc.UnaryServerInterceptor {
	return grpc_prometheus.UnaryServerInterceptor
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	head := metadata.Pairs("action", "123")
	_ = grpc.SetHeader(ctx, head)

	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)

	log.Printf("before invoker. method: %+v, request:%+v, md: %v", info.FullMethod, req, md)
	resp, err := handler(ctx, req)
	log.Printf("before invoker. method: %+v, response:%+v, duration: %v", info.FullMethod, resp, time.Since(start))

	return resp, err
}
