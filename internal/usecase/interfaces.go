package usecase

import (
	"context"
	"social/ent"
	"social/internal/entity"
)

type (
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
		GetGroupByIdRepo(ctx context.Context, groupId int64) (*ent.Group, error)
		CreateGroupRepo(ctx context.Context, group *entity.Group) (*ent.Group, error)
	}

	GroupMemberRepo interface {
		GetGroupsRepo(ctx context.Context, uid int64) ([]*entity.Group, error)
		GetGroupMembersRepo(ctx context.Context, uid int64) ([]*ent.GroupMember, error)
		CreateGroupMemberRepo(ctx context.Context, info *ent.GroupMember) (*ent.GroupMember, error)
		ExistsGroupMemberRepo(ctx context.Context, uid, groupId int64) (bool, error)
	}

	MessageRepo interface {
		AddMessageRepo(ctx context.Context, info *ent.Message) (*ent.Message, error)
		GetMessagesRepo(ctx context.Context, uid int64) ([]*ent.Message, error)
	}

	UserWebAPI interface{}
)
