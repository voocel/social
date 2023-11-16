package repo

import (
	"context"
	"fmt"
	"time"

	"social/ent"
	"social/ent/friendapply"
	"social/internal/entity"
	"social/pkg/log"
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
		SetCreatedAt(time.Now()).
		Save(ctx)
	return create, err
}

func (r FriendApplyRepo) GetFriendApplyRepo(ctx context.Context, uid int64) ([]*entity.FriendApplyResp, error) {
	found, err := r.ent.FriendApply.Query().Where(friendapply.ToID(uid)).Order(ent.Desc("id")).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetFriendApplyRepo query fail: %w", err)
	}
	res := make([]*entity.FriendApplyResp, 0)
	for _, apply := range found {
		user, err := r.ent.User.Get(ctx, apply.FromID)
		if err != nil {
			log.Errorf("GetFriendsRepo user get query fail: %w", err)
			continue
		}
		item := &entity.FriendApplyResp{
			FromId:     apply.FromID,
			FromName:   user.Nickname,
			FromAvatar: user.Avatar,
			ToId:       apply.ToID,
			Remark:     apply.Remark,
			Status:     apply.Status,
			CreatedAt:  apply.CreatedAt,
		}
		res = append(res, item)
	}
	return res, nil
}

func (r FriendApplyRepo) AgreeFriendApplyRepo(ctx context.Context, fromID, toID int64) (int, error) {
	found, err := r.ent.FriendApply.Update().
		Where(friendapply.FromID(fromID), friendapply.ToID(toID)).
		SetStatus(1).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("AgreeFriendApplyRepo update fail: %w", err)
	}

	return found, nil
}

func (r FriendApplyRepo) RefuseFriendApplyRepo(ctx context.Context, fromID, toID int64) (int, error) {
	found, err := r.ent.FriendApply.Update().
		Where(friendapply.FromID(fromID), friendapply.ToID(toID)).
		SetStatus(2).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("AgreeFriendApplyRepo update fail: %w", err)
	}
	return found, nil
}
