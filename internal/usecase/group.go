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

func (g *GroupUseCase) GetGroupsRepo(ctx context.Context, uid int64) ([]*ent.Group, error) {
	return g.repo.GetGroupsRepo(ctx, uid)
}

func (g *GroupUseCase) CreateGroupRepo(ctx context.Context, group *entity.Group) (*ent.Group, error) {
	return nil, nil
}
