package repo

import (
	"context"
	"fmt"
	"social/ent"
	entUser "social/ent/user"
	"social/internal/entity"
)

type UserRepo struct {
	ent *ent.Client
}

func NewUserRepo(ent *ent.Client) *UserRepo {
	return &UserRepo{ent}
}

func (r *UserRepo) GetUserByIdRepo(ctx context.Context, uid int64) (*ent.User, error) {
	found, err := r.ent.User.Query().Where(entUser.ID(uid)).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetUserByIdRepo query fail: %w", err)
	}
	return found, nil
}

func (r *UserRepo) GetUserByNameRepo(ctx context.Context, name string) (*ent.User, error) {
	found, err := r.ent.User.Query().Where(entUser.Username(name)).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetUserByNameRepo query fail: %w", err)
	}
	return found, err
}

func (r *UserRepo) GetUserByNameExistRepo(ctx context.Context, name string) (bool, error) {
	exist, err := r.ent.User.Query().Where(entUser.Username(name)).Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("UserRepo - GetUserByNameExistRepo query fail: %w", err)
	}
	return exist, err
}

func (r *UserRepo) GetUsersRepo(ctx context.Context) ([]*ent.User, int, error) {
	total, err := r.ent.User.Query().Where(entUser.DeletedAtIsNil()).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	users, err := r.ent.User.Query().
		Where(entUser.DeletedAtIsNil()).
		Limit(1).
		Offset(0).
		Order().All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepo) AddUserRepo(ctx context.Context, user *entity.User) (*ent.User, error) {
	create, err := r.ent.User.Create().
		SetUsername(user.Username).
		SetPassword(user.Password).
		SetMobile(user.Mobile).
		SetNickname(user.Nickname).
		SetEmail("").
		Save(ctx)
	return create, err
}

func (r *UserRepo) UpdateAvatarUserRepo(ctx context.Context, uid int64, avatar string) (int, error) {
	n, err := r.ent.User.Update().
		Where(entUser.ID(uid)).
		SetAvatar(avatar).
		Save(ctx)
	return n, err
}
