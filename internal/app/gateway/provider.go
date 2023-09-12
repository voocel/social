package gateway

import (
	"encoding/json"
	"social/internal/app/gateway/packet"
	"social/internal/entity"
	"social/internal/route"
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
		Seq:    0,
		Route:  route.Message,
		Buffer: req.Message.GetBuffer(),
	}
	b, _ := packet.Pack(data)
	return s.Push(b)
}

func (p provider) Broadcast(msg []byte) (n int64) {
	return p.gate.sessionGroup.Broadcast(msg)
}
