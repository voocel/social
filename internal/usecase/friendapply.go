package usecase

import (
	"context"
	"social/ent"
	"social/internal/entity"
)

type FriendApplyUseCase struct {
	repo FriendApplyRepo
}

func NewFriendApplyUseCase(fa FriendApplyRepo) *FriendApplyUseCase {
	return &FriendApplyUseCase{repo: fa}
}

func (f *FriendApplyUseCase) AddFriendApply(ctx context.Context, req *entity.FriendApply) (*ent.FriendApply, error) {
	return f.repo.AddFriendApplyRepo(ctx, req)
}

func (f *FriendApplyUseCase) GetFriendApply(ctx context.Context, uid int64) ([]*ent.FriendApply, error) {
	return f.repo.GetFriendApplyRepo(ctx, uid)
}

func (f *FriendApplyUseCase) AgreeFriendApply(ctx context.Context, applyId int64) error {
	return f.repo.AgreeFriendApplyRepo(ctx, applyId)
}
