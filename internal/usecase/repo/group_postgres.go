package repo

import (
	"context"
	"fmt"
	"social/ent"
	entGroup "social/ent/group"
	"social/internal/entity"
)

type GroupRepo struct {
	ent *ent.Client
}

func NewGroupRepo(ent *ent.Client) *GroupRepo {
	return &GroupRepo{ent: ent}
}

func (r *GroupRepo) GetGroupsRepo(ctx context.Context, uid int64) ([]*ent.Group, error) {
	found, err := r.ent.Group.Query().Where(entGroup.ID(uid)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("FriendRepo - GetFriendsRepo query fail: %w", err)
	}
	return found, nil
}

func (r *GroupRepo) CreateGroupRepo(ctx context.Context, group *entity.Group) (*ent.Group, error) {
	return nil, nil
}
