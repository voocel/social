package gate

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"social/internal/code"
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
func (gs *gateService) Bind(ctx context.Context, req *pb.BindRequest) (*pb.BindReply, error) {
	if req.Cid <= 0 || req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	s, err := gs.provider.Session(req.Uid)
	if err != nil {
		return nil, err
	}

	s.Bind(req.GetUid())

	return &pb.BindReply{}, nil
}

func (gs *gateService) Unbind(ctx context.Context, req *pb.UnbindRequest) (*pb.UnbindReply, error) {
	if req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	s, err := gs.provider.Session(req.Uid)
	if err != nil {
		return nil, err
	}
	s.Unbind(req.Uid)

	return &pb.UnbindReply{}, nil
}

// Push gateway send message to user
func (gs *gateService) Push(ctx context.Context, req *pb.PushRequest) (*pb.PushReply, error) {
	log.Debugf("[Gateway] receive node grpc message to user[%v]: %v", req.Target, string(req.GetMessage().GetBuffer()))
	resp := new(entity.Response)
	msg := &entity.Message{
		ID:          0,
		Content:     string(req.GetMessage().GetBuffer()),
		MsgType:     0,
		ContentType: 0,
	}
	s, err := gs.provider.Session(req.Target)
	if err != nil && errors.Is(err, code.ErrSessionNotFound) {
		// user offline
		entity.MsgCache.Add(req.Target, msg)

		st := status.New(codes.ResourceExhausted, "session does not exist")
		details, e := st.WithDetails(
			&errdetails.QuotaFailure{
				Violations: []*errdetails.QuotaFailure_Violation{{
					Subject:     strconv.Itoa(int(req.Target)),
					Description: code.ErrSessionNotFound.Error(),
				}},
			},
		)
		if e == nil {
			return &pb.PushReply{}, details.Err()
		}
		return &pb.PushReply{}, st.Err()
	}
	err = s.Push(resp.Resp(msg))
	if err != nil {
		log.Errorf("[Gateway] push to user(%v) err: ", req.Target, err)
	}
	return &pb.PushReply{}, nil
}

func (gs *gateService) GetIP(ctx context.Context, req *pb.GetIPRequest) (*pb.GetIPReply, error) {
	s, err := gs.provider.Session(req.Uid)
	if err != nil {
		return &pb.GetIPReply{}, nil
	}
	return &pb.GetIPReply{IP: s.RemoteIP()}, nil
}
