package usecase

import (
	"context"
	"social/ent"
)

type GroupMemberUseCase struct {
	repo GroupMemberRepo
}

func NewGroupMemberUseCase(r GroupMemberRepo) *GroupMemberUseCase {
	return &GroupMemberUseCase{repo: r}
}

func (g *GroupMemberUseCase) GetGroups(ctx context.Context, uid int64) ([]*ent.Group, error) {
	return g.repo.GetGroupsRepo(ctx, uid)
}

func (g *GroupMemberUseCase) GetGroupMemberUser(ctx context.Context, id int64) ([]int64, error) {
	return g.repo.GetGroupMemberUserRepo(ctx, id)
}

func (g *GroupMemberUseCase) CreateGroupMember(ctx context.Context, group *ent.GroupMember) (*ent.GroupMember, error) {
	return g.repo.CreateGroupMemberRepo(ctx, group)
}

func (g *GroupMemberUseCase) ExistsGroupMember(ctx context.Context, uid, groupId int64) (bool, error) {
	return g.repo.ExistsGroupMemberRepo(ctx, uid, groupId)
}
