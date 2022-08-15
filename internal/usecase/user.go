package usecase

import (
	"context"
	"fmt"
	"social/internal/entity"
)

type UserUseCase struct {
	repo   UserRepo
}

func NewUserUseCase(r UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) GetUserById(ctx context.Context, uid int) (*entity.User, error) {
	userInfo, err := u.repo.GetUserByIdRepo(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("GetUserByIdRepo err: %w", err)
	}
	return userInfo, nil
}