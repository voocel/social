package node

import (
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

func (p provider) Deliver(cid, uid int64, message *transport.Message) {
	route, ok := p.node.Routes[message.Route]
	r := Request{
		Cid:    cid,
		Uid:    uid,
		Seq:    message.Seq,
		Route:  message.Route,
		Buffer: message.Buffer,
		Node:   p.node,
	}

	if ok {
		route.Handler(r)
	} else if p.node.DefaultRouteHandler != nil {
		p.node.DefaultRouteHandler(r)
	} else {
		log.Errorf("[%v]the route does not match: %v", uid, message.Route)
	}
}
