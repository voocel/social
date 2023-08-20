package gate

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"social/internal/entity"
	"social/internal/transport"
	"social/pkg/log"
	"social/protos/pb"
)

type gateService struct {
	provider transport.GateProvider
	pb.UnimplementedGateServer
}

// Bind 将用户与当前网关进行绑定
func (e *gateService) Bind(ctx context.Context, req *pb.BindRequest) (*pb.BindReply, error) {
	if req.Cid <= 0 || req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	s, err := e.provider.Session(req.Uid)
	if err != nil {
		return nil, err
	}

	s.Bind(req.GetUid())

	return &pb.BindReply{}, nil
}

func (e *gateService) Unbind(ctx context.Context, req *pb.UnbindRequest) (*pb.UnbindReply, error) {
	if req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	s, err := e.provider.Session(req.Uid)
	if err != nil {
		return nil, err
	}
	s.Unbind(req.Uid)

	return &pb.UnbindReply{}, nil
}

// Push send to user
func (e *gateService) Push(ctx context.Context, req *pb.PushRequest) (*pb.PushReply, error) {
	log.Debugf("[Gateway] receive node grpc message to user[%v]: %v", req.Target, string(req.GetMessage().GetBuffer()))
	resp := new(entity.Response)
	msg := entity.Message{
		ID:          0,
		Content:     string(req.GetMessage().GetBuffer()),
		MsgType:     0,
		ContentType: 0,
	}
	s, err := e.provider.Session(req.Target)
	if err != nil {
		return nil, err
	}
	err = s.Push(resp.Resp(msg))
	if err != nil {
		log.Errorf("[Gateway] push to user(%v) err: ", req.Target, err)
	}
	return &pb.PushReply{}, nil
}
