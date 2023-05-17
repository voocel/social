package repo

import (
	"context"
	"fmt"
	"social/ent"
	entFriend "social/ent/friend"
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
