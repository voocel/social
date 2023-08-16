package node

import (
	"context"
	"fmt"
	node2 "social/internal/node"
	"social/pkg/log"
	"social/protos/pb"
)

// NodeService implemented grpc node server
type NodeService struct {
	pb.UnimplementedNodeServer
	Node *node2.Node
}

// Trigger Events triggered from the gateway
func (n NodeService) Trigger(ctx context.Context, req *pb.TriggerRequest) (*pb.TriggerReply, error) {
	n.Node.TriggerEvent(node2.Event(req.Event), req.GetGid(), req.GetUid())

	return &pb.TriggerReply{}, nil
}

// Deliver Messages sent from the gateway
func (n NodeService) Deliver(ctx context.Context, req *pb.DeliverRequest) (*pb.DeliverReply, error) {
	n.Node.RLock()
	route, ok := n.Node.Routes[req.Message.Route]
	n.Node.RUnlock()
	r := node2.Request{
		Gid:    req.Gid,
		Nid:    req.Nid,
		Cid:    req.Cid,
		Uid:    req.Uid,
		Route:  req.Message.Route,
		Buffer: req.Message.Buffer,
		Node:   n.Node,
	}
	var err error
	if ok {
		err = route.Handler(r)
	} else if n.Node.DefaultRouteHandler != nil {
		err = n.Node.DefaultRouteHandler(r)
	} else {
		err = fmt.Errorf("the route does not match: %v", req.Message.Route)
	}
	log.Debugf("[node] handle route err: %v", err)
	return &pb.DeliverReply{}, err
}
