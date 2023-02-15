package node

import (
	"context"
	"fmt"
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
	var err error
	if ok {
		err = route.handler(r)
	} else if n.node.defaultRouteHandler != nil {
		err = n.node.defaultRouteHandler(r)
	} else {
		err = fmt.Errorf("the route does not match: %v", req.Route)
	}
	log.Debugf("[node] handle route err: %v", err)
	return &node.DeliverReply{}, err
}
