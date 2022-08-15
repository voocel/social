package repo

import (
	"context"
	"social/internal/entity"
	"social/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) GetUserByIdRepo(ctx context.Context, uid int) (*entity.User, error) {
	return &entity.User{}, nil
}
