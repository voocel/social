package repo

import (
	"context"
	"fmt"
	"social/ent"
	"social/ent/groupmember"
	"social/internal/entity"
	"social/pkg/log"
)

type GroupMemberRepo struct {
	ent *ent.Client
}

func NewGroupMemberRepo(ent *ent.Client) *GroupMemberRepo {
	return &GroupMemberRepo{ent: ent}
}

func (g *GroupMemberRepo) GetGroupsRepo(ctx context.Context, uid int64) ([]*entity.Group, error) {
	found, err := g.ent.GroupMember.Query().Where(groupmember.UID(uid)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetGroupMembersRepo query fail: %w", err)
	}
	res := make([]*entity.Group, 0)
	for _, member := range found {
		group, err := g.ent.Group.Get(ctx, member.GroupID)
		if err != nil {
			log.Errorf("GetGroupMembersRepo user get query fail: %w", err)
			continue
		}
		item := &entity.Group{
			Name:         group.Name,
			Owner:        group.Owner,
			Notice:       group.Notice,
			Introduction: group.Introduction,
			CreatedUid:   group.CreatedUID,
		}
		res = append(res, item)
	}
	return res, nil
}

func (g *GroupMemberRepo) GetGroupMembersRepo(ctx context.Context, uid int64) ([]*ent.GroupMember, error) {
	return g.ent.GroupMember.Query().Where(groupmember.UID(uid)).All(ctx)
}

func (g *GroupMemberRepo) CreateGroupMemberRepo(ctx context.Context, info *ent.GroupMember) (*ent.GroupMember, error) {
	return g.ent.GroupMember.Create().
		SetGroupID(info.GroupID).
		SetUID(info.UID).
		Save(ctx)
}
