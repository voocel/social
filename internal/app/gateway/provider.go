package gateway

import (
	"social/internal/entity"
	"social/internal/session"
)

type provider struct {
	gate *Gateway
}

func (p provider) Session(target int64) (*session.Session, error) {
	return p.gate.sessionGroup.GetSessionByUid(target)
}

func (p provider) Push(target int64, msg []byte) error {
	data := &entity.Message{
		ID:          0,
		Content:     string(msg),
		MsgType:     0,
		ContentType: 0,
	}
	s, err := p.gate.sessionGroup.GetSessionByUid(target)
	if err != nil {
		// user offline
		entity.MsgCache.Add(target, data)
		return err
	}
	resp := new(entity.Response)
	return s.Push(resp.Resp(data))
}

func (p provider) Broadcast(msg []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}
