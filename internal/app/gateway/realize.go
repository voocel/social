package gateway

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"social/internal/entity"
	"social/pkg/log"
	"social/protos/gate"
)

type endpoint struct {
	sessionGroup *SessionGroup
	gate.UnimplementedGateServer
}

// Bind 将用户与当前网关进行绑定
func (e endpoint) Bind(ctx context.Context, req *gate.BindRequest) (*gate.BindReply, error) {
	if req.Cid <= 0 || req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	return &gate.BindReply{}, nil
}

func (e endpoint) Unbind(ctx context.Context, req *gate.UnbindRequest) (*gate.UnbindReply, error) {
	panic("implement me")
}

// Push send to user
func (e endpoint) Push(ctx context.Context, req *gate.PushRequest) (*gate.PushReply, error) {
	log.Debugf("[Gateway] receive node grpc message to user[%v]: %v", req.Target, string(req.GetBuffer()))
	resp := new(entity.Response)
	msg := entity.Message{
		ID:          0,
		Content:     string(req.GetBuffer()),
		MsgType:     0,
		ContentType: 0,
	}
	err := e.sessionGroup.PushByUid(req.Target, resp.Resp(msg))
	if err != nil {
		log.Errorf("[Gateway] push to user(%v) err: ", req.Target, err)
	}
	return &gate.PushReply{}, nil
}
