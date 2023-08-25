package gateway

import (
	"context"
	"social/internal/session"
)

type provider struct {
	gate *Gateway
}

func (p provider) Session(target int64) (*session.Session, error) {
	return p.gate.sessionGroup.GetSessionByUid(target)
}

func (p provider) Bind(ctx context.Context, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (p provider) Unbind(ctx context.Context, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (p provider) Push(target int64, msg []byte, msgType ...int) error {
	//TODO implement me
	panic("implement me")
}

func (p provider) Broadcast(msg []byte, msgType ...int) (n int, err error) {
	//TODO implement me
	panic("implement me")
}
