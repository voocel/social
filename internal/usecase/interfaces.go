package usecase

import (
	"context"
	"social/ent"
	"social/internal/entity"
)

type (
	User interface {
		GetUserById(ctx context.Context, uid int64) (*ent.User, error)
	}

	UserRepo interface {
		GetUserByIdRepo(ctx context.Context, uid int64) (*ent.User, error)
		GetUserByNameRepo(ctx context.Context, name string) (*ent.User, error)
		GetUserByNameExistRepo(ctx context.Context, name string) (bool, error)
		GetUsersRepo(ctx context.Context) ([]*ent.User, int, error)
		AddUserRepo(ctx context.Context, user *entity.User) (*ent.User, error)
		UpdateAvatarUserRepo(ctx context.Context, uid int64, avatar string) (int, error)
	}

	FriendRepo interface {
		GetFriendsRepo(ctx context.Context, uid int64) ([]*entity.FriendResp, error)
		AddFriendRepo(ctx context.Context, friend *entity.Friend) (*ent.Friend, error)
		DeleteFriendRepo(ctx context.Context, uid, friendId int64) error
	}

	FriendApplyRepo interface {
		AddFriendApplyRepo(ctx context.Context, req *entity.FriendApply) (*ent.FriendApply, error)
		GetFriendApplyRepo(ctx context.Context, uid int64) ([]*entity.FriendApplyResp, error)
		AgreeFriendApplyRepo(ctx context.Context, fromID, toID int64) (int, error)
		RefuseFriendApplyRepo(ctx context.Context, fromID, toID int64) (int, error)
	}

	GroupRepo interface {
		GetGroupsRepo(ctx context.Context, uid int64) ([]*ent.Group, error)
		CreateGroupRepo(ctx context.Context, group *entity.Group) (*ent.Group, error)
	}

	UserWebAPI interface{}
)
