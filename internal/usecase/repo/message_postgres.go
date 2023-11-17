package repo

import (
	"context"

	"social/ent"
	"social/ent/message"
)

type MessageRepo struct {
	ent *ent.Client
}

func NewMessageRepo(ent *ent.Client) *MessageRepo {
	return &MessageRepo{ent}
}

func (m *MessageRepo) AddMessageRepo(ctx context.Context, info *ent.Message) (*ent.Message, error) {
	return m.ent.Message.Create().
		SetSenderID(info.SenderID).
		SetReceiverID(info.ReceiverID).
		SetContent(info.Content).
		Save(ctx)
}

func (m *MessageRepo) GetMessagesRepo(ctx context.Context, uid int64) ([]*ent.Message, error) {
	return m.ent.Message.Query().Where(message.ReceiverID(uid)).All(ctx)
}
