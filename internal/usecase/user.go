package usecase

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"social/ent"
	"social/internal/entity"
	"social/pkg/util"
	"time"
)

type UserUseCase struct {
	repo UserRepo
	sf   singleflight.Group
}

func NewUserUseCase(r UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) UserLogin(ctx context.Context, req entity.UserLoginReq) (*ent.User, error) {
	user, err := u.repo.GetUserByNameRepo(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user.Status != 1 {
		return nil, errors.New("this user is disable")
	}
	if !util.VerifyPassword(req.Password, user.Password) {
		return nil, errors.New("password incorrect")
	}
	return user, nil
}

func (u *UserUseCase) UserRegister(ctx context.Context, req entity.UserRegisterReq) error {
	exist, err := u.repo.GetUserByNameExistRepo(ctx, req.Username)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("username already exists")
	}
	userInfo := &entity.User{}
	userInfo.Password, err = util.EncryptPassword(req.Password)
	if err != nil {
		return err
	}
	userInfo.Username = req.Username
	userInfo.Nickname = "unknown"
	userInfo.LastLoginTime = time.Now()
	_, err = u.repo.AddUserRepo(ctx, userInfo)
	return err
}

func (u *UserUseCase) GetUserById(ctx context.Context, uid int64) (*ent.User, error) {
	v, err, _ := u.sf.Do("key1", func() (interface{}, error) {
		userInfo, err := u.repo.GetUserByIdRepo(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("GetUserByIdRepo err: %w", err)
		}
		// todo set cache
		return userInfo, err
	})

	return v.(*ent.User), err
}

func (u *UserUseCase) GetUserByName(ctx context.Context, name string) (*ent.User, error) {
	v, err, _ := u.sf.Do(name, func() (interface{}, error) {
		userInfo, err := u.repo.GetUserByNameRepo(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("GetUserByIdRepo err: %w", err)
		}
		// todo set cache
		return userInfo, err
	})

	return v.(*ent.User), err
}
