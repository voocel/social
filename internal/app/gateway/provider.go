package gateway

import (
	"encoding/json"
	"social/internal/app/gateway/packet"
	"social/internal/entity"
	"social/internal/session"
	"social/protos/pb"
)

type provider struct {
	gate *Gateway
}

func (p provider) Session(target int64) (*session.Session, error) {
	return p.gate.sessionGroup.GetSessionByUid(target)
}

// Push sent message to user client
func (p provider) Push(req *pb.PushReq) error {
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(req.Message.GetBuffer(), msg); err != nil {
		return err
	}
	s, err := p.gate.sessionGroup.GetSessionByUid(req.GetTarget())
	if err != nil {
		// user offline
		entity.MsgCache.Add(req.GetTarget(), msg)
		return err
	}

	data := &packet.Message{
		Seq:    req.Message.Seq,
		Route:  req.Message.Route,
		Buffer: req.Message.GetBuffer(),
	}
	b, _ := packet.Pack(data)
	return s.Push(b)
}

func (p provider) Broadcast(req *pb.Message) (n int64) {
	data := &packet.Message{
		Seq:    req.Seq,
		Route:  req.Route,
		Buffer: req.GetBuffer(),
	}
	b, _ := packet.Pack(data)
	return p.gate.sessionGroup.Broadcast(b)
}
