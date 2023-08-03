package repo

import (
	"context"
	"fmt"
	"social/ent"
	entFriendApply "social/ent/friendapply"
	"social/internal/entity"
)

type FriendApplyRepo struct {
	ent *ent.Client
}

func NewFriendApplyRepo(ent *ent.Client) *FriendApplyRepo {
	return &FriendApplyRepo{ent}
}

func (r FriendApplyRepo) AddFriendApplyRepo(ctx context.Context, req *entity.FriendApply) (*ent.FriendApply, error) {
	create, err := r.ent.FriendApply.Create().
		SetFromID(req.FromId).
		SetToID(req.ToId).
		SetRemark(req.Remark).
		SetStatus(0).
		Save(ctx)
	return create, err
}

func (r FriendApplyRepo) GetFriendApplyRepo(ctx context.Context, uid int64) ([]*ent.FriendApply, error) {
	found, err := r.ent.FriendApply.Query().Where(entFriendApply.ToID(uid)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("FriendApplyRepo - GetFriendApplyRepo query fail: %w", err)
	}
	return found, nil
}

func (r FriendApplyRepo) AgreeFriendApplyRepo(ctx context.Context, fromID, toID int64) (*ent.FriendApply, error) {
	found, err := r.ent.FriendApply.UpdateOne(&ent.FriendApply{
		FromID: fromID,
		ToID:   toID,
	}).SetStatus(1).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("FriendApplyRepo - AgreeFriendApplyRepo update fail: %w", err)
	}
	return found, nil
}

func (r FriendApplyRepo) RefuseFriendApplyRepo(ctx context.Context, applyId int64) error {
	panic("implement me")
}
