package node

import (
	"fmt"
	"social/internal/event"
	"social/internal/transport"
	"social/pkg/log"
)

type provider struct {
	node *Node
}

func (p provider) Trigger(event event.Event, gid string, uid int64) {
	p.node.triggerEvent(event, gid, uid)
}

func (p provider) Deliver(gid, nid string, cid, uid int64, message *transport.Message) {
	route, ok := p.node.Routes[message.Route]
	r := Request{
		Gid:    gid,
		Nid:    nid,
		Cid:    cid,
		Uid:    uid,
		Route:  message.Route,
		Buffer: message.Buffer,
		Node:   p.node,
	}
	var err error
	if ok {
		err = route.Handler(r)
	} else if p.node.DefaultRouteHandler != nil {
		err = p.node.DefaultRouteHandler(r)
	} else {
		err = fmt.Errorf("the route does not match: %v", message.Route)
	}
	log.Debugf("[node] handle route err: %v", err)
}
