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
	}

	FriendRepo interface {
		GetFriendsRepo(ctx context.Context, uid int64) ([]*ent.Friend, error)
		AddFriendRepo(ctx context.Context, friend *entity.Friend) (*ent.Friend, error)
	}

	UserWebAPI interface{}
)
