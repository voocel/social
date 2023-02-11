package node

import (
	"context"
	"social/pkg/log"
	"social/protos/node"
)

// implemented grpc node server
type nodeService struct {
	node.UnimplementedNodeServer
	node *Node
}

// Trigger Events triggered from the gateway
func (n nodeService) Trigger(ctx context.Context, req *node.TriggerRequest) (*node.TriggerReply, error) {
	n.node.triggerEvent(Event(req.Event), req.GetGid(), req.GetUid())

	return &node.TriggerReply{}, nil
}

// Deliver Messages sent from the gateway
func (n nodeService) Deliver(ctx context.Context, req *node.DeliverRequest) (*node.DeliverReply, error) {
	n.node.RLock()
	route, ok := n.node.routes[req.Route]
	n.node.RUnlock()
	r := Request{
		Gid:    req.Gid,
		Nid:    req.Nid,
		Cid:    req.Cid,
		Uid:    req.Uid,
		Route:  req.Route,
		Buffer: req.Buffer,
		Node:   n.node,
	}
	if ok {
		route.handler(r)
	} else if n.node.defaultRouteHandler != nil {
		n.node.defaultRouteHandler(r)
	} else {
		log.Errorf("[node] the route does not match: %v", req.Route)
	}

	return &node.DeliverReply{}, nil
}
