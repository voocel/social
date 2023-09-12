package usecase

import (
	"context"
	"golang.org/x/sync/singleflight"
	"social/ent"
	"social/internal/entity"
)

type FriendUseCase struct {
	repo FriendRepo
	sf   singleflight.Group
}

func NewFriendUseCase(r FriendRepo) *FriendUseCase {
	return &FriendUseCase{repo: r}
}

func (f *FriendUseCase) GetFriends(ctx context.Context, uid int64) ([]*entity.FriendResp, error) {
	v, err, _ := f.sf.Do("key", func() (interface{}, error) {
		return f.repo.GetFriendsRepo(ctx, uid)
	})
	return v.([]*entity.FriendResp), err
}

func (f *FriendUseCase) AddFriend(ctx context.Context, req *entity.Friend) (*ent.Friend, error) {
	return f.repo.AddFriendRepo(ctx, req)
}
