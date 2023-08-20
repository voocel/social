package node

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"social/internal/node"
	"social/internal/transport"
	"social/protos/pb"
)

// NodeService implemented grpc node server
type nodeService struct {
	pb.UnimplementedNodeServer
	provider transport.NodeProvider
}

// Trigger Events triggered from the gateway
func (n nodeService) Trigger(ctx context.Context, req *pb.TriggerRequest) (*pb.TriggerReply, error) {
	if req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}

	n.provider.Trigger(node.Event(req.Event), req.Gid, req.Uid)

	return &pb.TriggerReply{}, nil
}

// Deliver Messages sent from the gateway
func (n nodeService) Deliver(ctx context.Context, req *pb.DeliverRequest) (*pb.DeliverReply, error) {
	if req.Uid <= 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument").Err()
	}
	n.provider.Deliver(req.Gid, req.Nid, req.Cid, req.Uid, &transport.Message{
		Seq:    req.Message.Seq,
		Route:  req.Message.Route,
		Buffer: req.Message.Buffer,
	})

	return &pb.DeliverReply{}, nil
}
