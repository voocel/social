package usecase

import (
	"context"
	"social/ent"
	"social/internal/entity"
)

type (
	User interface {
		GetUserById(ctx context.Context, uid int64) (*entity.User, error)
	}

	UserRepo interface {
		GetUserByIdRepo(ctx context.Context, uid int64) (*ent.User, error)
		GetUserByNameRepo(ctx context.Context, name string) (*ent.User, error)
		GetUsersRepo(ctx context.Context) ([]*ent.User, int, error)
	}

	UserWebAPI interface{}
)
