package repo

import (
	"context"
	"fmt"

	"social/ent"
	"social/ent/groupmember"
	"social/pkg/log"
)

type GroupMemberRepo struct {
	ent *ent.Client
}

func NewGroupMemberRepo(ent *ent.Client) *GroupMemberRepo {
	return &GroupMemberRepo{ent: ent}
}

func (g *GroupMemberRepo) GetGroupsRepo(ctx context.Context, uid int64) ([]*ent.Group, error) {
	found, err := g.ent.GroupMember.Query().Where(groupmember.UID(uid)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetGroupMembersRepo query fail: %w", err)
	}
	res := make([]*ent.Group, 0)
	for _, member := range found {
		group, err := g.ent.Group.Get(ctx, member.GroupID)
		if err != nil {
			log.Errorf("GetGroupMembersRepo user get query fail: %w", err)
			continue
		}
		res = append(res, group)
	}
	return res, nil
}

func (g *GroupMemberRepo) GetGroupMembersRepo(ctx context.Context, uid int64) ([]*ent.GroupMember, error) {
	return g.ent.GroupMember.Query().Where(groupmember.UID(uid)).All(ctx)
}

func (g *GroupMemberRepo) GetGroupMemberUserRepo(ctx context.Context, id int64) ([]int64, error) {
	found, err := g.ent.GroupMember.Query().Select(groupmember.FieldUID).Where(groupmember.GroupID(id)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetGroupMemberUserRepo query fail: %w", err)
	}
	res := make([]int64, 0)
	for _, item := range found {
		res = append(res, item.UID)
	}
	return res, nil
}

func (g *GroupMemberRepo) CreateGroupMemberRepo(ctx context.Context, info *ent.GroupMember) (*ent.GroupMember, error) {
	return g.ent.GroupMember.Create().
		SetGroupID(info.GroupID).
		SetUID(info.UID).
		SetRemark(info.Remark).
		Save(ctx)
}

func (g *GroupMemberRepo) ExistsGroupMemberRepo(ctx context.Context, uid, groupId int64) (bool, error) {
	return g.ent.GroupMember.Query().Where(groupmember.UID(uid), groupmember.GroupID(groupId)).Exist(ctx)
}
