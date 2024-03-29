package gate

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"social/internal/code"
	"social/internal/transport"
	"social/pkg/log"
	"social/protos/pb"
)

type gateService struct {
	provider transport.GateProvider
	pb.UnimplementedGateServer
}

// Bind 将用户与当前网关进行绑定
func (gs *gateService) Bind(ctx context.Context, req *pb.BindReq) (*pb.BaseReply, error) {
	if req.Cid <= 0 || req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	s, err := gs.provider.Session(req.Uid)
	if err != nil {
		return nil, err
	}

	s.Bind(req.GetUid())

	return &pb.BaseReply{}, nil
}

func (gs *gateService) Unbind(ctx context.Context, req *pb.BaseReq) (*pb.BaseReply, error) {
	if req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	s, err := gs.provider.Session(req.Uid)
	if err != nil {
		return nil, err
	}
	s.Unbind(req.Uid)

	return &pb.BaseReply{}, nil
}

func (gs *gateService) GetIP(ctx context.Context, req *pb.BaseReq) (*pb.GetIPReply, error) {
	s, err := gs.provider.Session(req.Uid)
	if err != nil {
		return &pb.GetIPReply{}, nil
	}
	return &pb.GetIPReply{Ip: s.RemoteIP()}, nil
}

// Push gateway send message to user
func (gs *gateService) Push(ctx context.Context, req *pb.PushReq) (*pb.BaseReply, error) {
	log.Debugf("[Gateway] receive node grpc message to user[%v]: %v", req.Target, string(req.GetMessage().GetBuffer()))
	err := gs.provider.Push(req)
	if err != nil && errors.Is(err, code.ErrSessionNotFound) {
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
			return &pb.BaseReply{}, details.Err()
		}
		return &pb.BaseReply{}, st.Err()
	}
	if err != nil {
		log.Errorf("[Gateway] push to user(%v) err: ", req.Target, err)
	}
	return &pb.BaseReply{}, nil
}

func (gs *gateService) Multicast(ctx context.Context, req *pb.MulticastReq) (*pb.MulticastReply, error) {
	log.Debugf("[Gateway] receive node grpc multicast message to target[%v]: %v", req.Targets, string(req.GetMessage().GetBuffer()))
	n := gs.provider.Multicast(req.Targets, req.Message)
	return &pb.MulticastReply{Total: n}, nil
}

func (gs *gateService) Broadcast(ctx context.Context, req *pb.BroadcastReq) (*pb.BroadcastReply, error) {
	n := gs.provider.Broadcast(req.Message)
	return &pb.BroadcastReply{Total: n}, nil
}

func (gs *gateService) Disconnect(ctx context.Context, req *pb.BaseReq) (*pb.BaseReply, error) {
	if req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	s, err := gs.provider.Session(req.Uid)
	if err != nil {
		return nil, err
	}
	err = s.Close()
	return &pb.BaseReply{}, err
}
