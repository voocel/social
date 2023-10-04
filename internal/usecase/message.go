package usecase

import (
	"context"
	"social/ent"
)

type MessageUseCase struct {
	repo MessageRepo
}

func NewMessageUseCase(m MessageRepo) *MessageUseCase {
	return &MessageUseCase{repo: m}
}

func (m *MessageUseCase) AddMessage(ctx context.Context, info *ent.Message) (*ent.Message, error) {
	return m.repo.AddMessageRepo(ctx, info)
}

func (m *MessageUseCase) GetMessages(ctx context.Context, uid int64) ([]*ent.Message, error) {
	return m.repo.GetMessagesRepo(ctx, uid)
}
