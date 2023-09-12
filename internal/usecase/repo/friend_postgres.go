package repo

import (
	"context"
	"fmt"
	"social/ent"
	entFriend "social/ent/friend"
	"social/internal/entity"
	"social/pkg/log"
)

type FriendRepo struct {
	ent *ent.Client
}

func NewFriendRepo(ent *ent.Client) *FriendRepo {
	return &FriendRepo{ent}
}

func (r *FriendRepo) GetFriendsRepo(ctx context.Context, uid int64) ([]*entity.FriendResp, error) {
	found, err := r.ent.Friend.Query().Where(entFriend.UID(uid)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetFriendsRepo query fail: %w", err)
	}

	res := make([]*entity.FriendResp, 0)
	for _, friend := range found {
		user, err := r.ent.User.Get(ctx, friend.FriendID)
		if err != nil {
			log.Errorf("GetFriendsRepo user get query fail: %w", err)
			continue
		}
		item := &entity.FriendResp{
			Uid:      uid,
			FriendId: friend.FriendID,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Remark:   friend.Remark,
			Shield:   friend.Shield,
		}
		res = append(res, item)
	}
	return res, nil
}

func (r *FriendRepo) AddFriendRepo(ctx context.Context, info *entity.Friend) (*ent.Friend, error) {
	tx, err := r.ent.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("AddFriendRepo start transaction fail: %w", err)
	}

	create, err := tx.Friend.Create().
		SetUID(info.Uid).
		SetFriendID(info.FriendId).
		SetRemark(info.Remark).
		Save(ctx)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			err = fmt.Errorf("%v: %v", err, e)
		}
		return nil, fmt.Errorf("AddFriendRepo create fail: %w", err)
	}

	_, err = tx.Friend.Create().
		SetUID(info.FriendId).
		SetFriendID(info.Uid).
		SetRemark(info.Remark).
		Save(ctx)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			err = fmt.Errorf("%v: %v", err, e)
		}
		return nil, fmt.Errorf("AddFriendRepo create2 fail: %w", err)
	}

	return create, tx.Commit()
}
