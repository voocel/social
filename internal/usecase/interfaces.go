package usecase

import (
	"context"
	"social/internal/entity"
)

type (
	User interface {
		GetUserById(ctx context.Context, uid int)(*entity.User, error)
	}

	UserRepo interface {
		GetUserByIdRepo(ctx context.Context, uid int)(*entity.User, error)
	}

	UserWebAPI interface {

	}
)
