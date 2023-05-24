package repo

import (
	"context"
	"fmt"
	"social/ent"
	entFriend "social/ent/friend"
	"social/internal/entity"
)

type FriendRepo struct {
	ent *ent.Client
}

func NewFriendRepo(ent *ent.Client) *FriendRepo {
	return &FriendRepo{ent}
}

func (r *FriendRepo) GetFriendsRepo(ctx context.Context, uid int64) ([]*ent.Friend, error) {
	found, err := r.ent.Friend.Query().Where(entFriend.ID(uid)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("FriendRepo - GetFriendsRepo query fail: %w", err)
	}
	return found, nil
}

func (r *FriendRepo) AddFriendRepo(ctx context.Context, info *entity.Friend) (*ent.Friend, error) {
	create, err := r.ent.Friend.Create().
		SetUID(info.Uid).
		SetFriendID(info.FriendId).
		SetRemark(info.Remark).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("FriendRepo - AddFriendRepo create fail: %w", err)
	}
	return create, nil
}
