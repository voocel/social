package usecase

import "context"

type FriendUseCase struct {
	repo FriendRepo
}

func NewFriendUseCase(r FriendRepo) *FriendUseCase {
	return &FriendUseCase{repo: r}
}

func (f *FriendUseCase) GetFriendsRepo(ctx context.Context, uid int64) {

}
