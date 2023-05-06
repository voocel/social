package usecase

import (
	"context"
	"social/ent"
)

type FriendUseCase struct {
	repo FriendRepo
}

func NewFriendUseCase(r FriendRepo) *FriendUseCase {
	return &FriendUseCase{repo: r}
}

func (f *FriendUseCase) GetFriendsRepo(ctx context.Context, uid int64) ([]*ent.Friend, error) {
	return f.repo.GetFriendsRepo(ctx, uid)
}
