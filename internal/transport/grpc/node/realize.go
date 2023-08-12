package node

import (
	"context"
	"fmt"
	node2 "social/internal/node"
	"social/pkg/log"
	"social/protos/node"
)

// NodeService implemented grpc node server
type NodeService struct {
	node.UnimplementedNodeServer
	Node *node2.Node
}

// Trigger Events triggered from the gateway
func (n NodeService) Trigger(ctx context.Context, req *node.TriggerRequest) (*node.TriggerReply, error) {
	n.Node.TriggerEvent(node2.Event(req.Event), req.GetGid(), req.GetUid())

	return &node.TriggerReply{}, nil
}

// Deliver Messages sent from the gateway
func (n NodeService) Deliver(ctx context.Context, req *node.DeliverRequest) (*node.DeliverReply, error) {
	n.Node.RLock()
	route, ok := n.Node.Routes[req.Route]
	n.Node.RUnlock()
	r := node2.Request{
		Gid:    req.Gid,
		Nid:    req.Nid,
		Cid:    req.Cid,
		Uid:    req.Uid,
		Route:  req.Route,
		Buffer: req.Buffer,
		Node:   n.Node,
	}
	var err error
	if ok {
		err = route.Handler(r)
	} else if n.Node.DefaultRouteHandler != nil {
		err = n.Node.DefaultRouteHandler(r)
	} else {
		err = fmt.Errorf("the route does not match: %v", req.Route)
	}
	log.Debugf("[node] handle route err: %v", err)
	return &node.DeliverReply{}, err
}
