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

func (f *FriendUseCase) GetFriendsRepo(ctx context.Context, uid int64) ([]*ent.Friend, error) {
	v, err, _ := f.sf.Do("key", func() (interface{}, error) {
		return f.repo.GetFriendsRepo(ctx, uid)
	})
	return v.([]*ent.Friend), err
}

func (f *FriendUseCase) AddFriend(ctx context.Context, req *entity.Friend) error {
	f.repo.AddFriendRepo(ctx, req)
	return nil
}
