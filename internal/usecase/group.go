package usecase

import (
	"context"
	"social/ent"
	"social/internal/entity"
)

type GroupUseCase struct {
}

func NewGroupUseCase() *GroupUseCase {
	return &GroupUseCase{}
}

func (g *GroupUseCase) GetGroupsRepo(ctx context.Context, uid int64) ([]*ent.Group, error) {
	return nil, nil

}

func (g *GroupUseCase) CreateGroupRepo(ctx context.Context, group *entity.Group) (*ent.Group, error) {
	return nil, nil
}
