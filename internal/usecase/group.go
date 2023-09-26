package usecase

import (
	"context"
	"social/ent"
	"social/internal/entity"
)

type GroupUseCase struct {
	repo GroupRepo
}

func NewGroupUseCase(r GroupRepo) *GroupUseCase {
	return &GroupUseCase{repo: r}
}

func (g *GroupUseCase) GetGroupById(ctx context.Context, groupId int64) (*ent.Group, error) {
	return g.repo.GetGroupByIdRepo(ctx, groupId)
}

func (g *GroupUseCase) GetGroups(ctx context.Context, uid int64) ([]*ent.Group, error) {
	return g.repo.GetGroupsRepo(ctx, uid)
}

func (g *GroupUseCase) CreateGroup(ctx context.Context, group *entity.Group) (*ent.Group, error) {
	return g.repo.CreateGroupRepo(ctx, group)
}
