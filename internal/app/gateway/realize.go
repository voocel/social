package gateway

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"social/pkg/log"
	"social/protos/gate"
)

type endpoint struct {
	sessionGroup *SessionGroup
	gate.UnimplementedGateServer
}

// Bind 将用户与当前网关进行绑定
func (e endpoint) Bind(ctx context.Context, req *gate.BindRequest) (*gate.BindReply, error) {
	if req.CID <= 0 || req.UID <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	return &gate.BindReply{}, nil
}

func (e endpoint) Unbind(ctx context.Context, req *gate.UnbindRequest) (*gate.UnbindReply, error) {
	panic("implement me")
}

func (e endpoint) Push(ctx context.Context, req *gate.PushRequest) (*gate.PushReply, error) {
	log.Debugf("[Gateway] receive node grpc message to user[%v]: %v", req.Target, string(req.GetBuffer()))
	err := e.sessionGroup.PushByUid(req.Target, req.GetBuffer())
	if err != nil {
		log.Errorf("[Gateway] push to user(%v) err: ", req.Target, err)
	}
	return &gate.PushReply{}, nil
}
