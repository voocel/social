package client

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	traceID := uuid.New().String()
	log.Printf("before invoker. method: %+v, request:%+v, trace_id: %v", method, req, traceID)
	var header, trailer metadata.MD
	opts = append(opts, grpc.Header(&header), grpc.Trailer(&trailer))
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("after invoker. reply: %+v, metadata: %+v", reply, header)
	return err
}

func UnaryClientRetryInterceptor() grpc.UnaryClientInterceptor {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(time.Second * 1),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinearWithJitter(500*time.Millisecond, 0.2)),
		grpc_retry.WithCodes(codes.Unavailable, codes.Aborted),
	}
	return grpc_retry.UnaryClientInterceptor(opts...)
}
