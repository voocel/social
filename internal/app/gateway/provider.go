package gateway

import (
	"encoding/json"
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
func (p provider) Push(target int64, buf []byte) error {
	var msg = new(pb.MsgItem)
	if err := json.Unmarshal(buf, msg); err != nil {
		return err
	}
	s, err := p.gate.sessionGroup.GetSessionByUid(target)
	if err != nil {
		// user offline
		entity.MsgCache.Add(target, msg)
		return err
	}
	resp := new(entity.Response)
	return s.Push(resp.Wrap(msg))
}

func (p provider) Broadcast(msg []byte) (n int) {
	return p.gate.sessionGroup.Broadcast(msg)
}
