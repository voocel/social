// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"social/ent/groupmember"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// GroupMemberCreate is the builder for creating a GroupMember entity.
type GroupMemberCreate struct {
	config
	mutation *GroupMemberMutation
	hooks    []Hook
}

// SetUID sets the "uid" field.
func (gmc *GroupMemberCreate) SetUID(i int64) *GroupMemberCreate {
	gmc.mutation.SetUID(i)
	return gmc
}

// SetGroupID sets the "group_id" field.
func (gmc *GroupMemberCreate) SetGroupID(i int64) *GroupMemberCreate {
	gmc.mutation.SetGroupID(i)
	return gmc
}

// SetInviter sets the "inviter" field.
func (gmc *GroupMemberCreate) SetInviter(i int64) *GroupMemberCreate {
	gmc.mutation.SetInviter(i)
	return gmc
}

// SetRemark sets the "remark" field.
func (gmc *GroupMemberCreate) SetRemark(s string) *GroupMemberCreate {
	gmc.mutation.SetRemark(s)
	return gmc
}

// SetStatus sets the "status" field.
func (gmc *GroupMemberCreate) SetStatus(i int8) *GroupMemberCreate {
	gmc.mutation.SetStatus(i)
	return gmc
}

// SetApplyAt sets the "apply_at" field.
func (gmc *GroupMemberCreate) SetApplyAt(t time.Time) *GroupMemberCreate {
	gmc.mutation.SetApplyAt(t)
	return gmc
}

// SetNillableApplyAt sets the "apply_at" field if the given value is not nil.
func (gmc *GroupMemberCreate) SetNillableApplyAt(t *time.Time) *GroupMemberCreate {
	if t != nil {
		gmc.SetApplyAt(*t)
	}
	return gmc
}

// SetCreatedAt sets the "created_at" field.
func (gmc *GroupMemberCreate) SetCreatedAt(t time.Time) *GroupMemberCreate {
	gmc.mutation.SetCreatedAt(t)
	return gmc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (gmc *GroupMemberCreate) SetNillableCreatedAt(t *time.Time) *GroupMemberCreate {
	if t != nil {
		gmc.SetCreatedAt(*t)
	}
	return gmc
}

// SetUpdatedAt sets the "updated_at" field.
func (gmc *GroupMemberCreate) SetUpdatedAt(t time.Time) *GroupMemberCreate {
	gmc.mutation.SetUpdatedAt(t)
	return gmc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (gmc *GroupMemberCreate) SetNillableUpdatedAt(t *time.Time) *GroupMemberCreate {
	if t != nil {
		gmc.SetUpdatedAt(*t)
	}
	return gmc
}

// SetDeletedAt sets the "deleted_at" field.
func (gmc *GroupMemberCreate) SetDeletedAt(t time.Time) *GroupMemberCreate {
	gmc.mutation.SetDeletedAt(t)
	return gmc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (gmc *GroupMemberCreate) SetNillableDeletedAt(t *time.Time) *GroupMemberCreate {
	if t != nil {
		gmc.SetDeletedAt(*t)
	}
	return gmc
}

// SetID sets the "id" field.
func (gmc *GroupMemberCreate) SetID(i int64) *GroupMemberCreate {
	gmc.mutation.SetID(i)
	return gmc
}

// Mutation returns the GroupMemberMutation object of the builder.
func (gmc *GroupMemberCreate) Mutation() *GroupMemberMutation {
	return gmc.mutation
}

// Save creates the GroupMember in the database.
func (gmc *GroupMemberCreate) Save(ctx context.Context) (*GroupMember, error) {
	var (
		err  error
		node *GroupMember
	)
	if len(gmc.hooks) == 0 {
		if err = gmc.check(); err != nil {
			return nil, err
		}
		node, err = gmc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*GroupMemberMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = gmc.check(); err != nil {
				return nil, err
			}
			gmc.mutation = mutation
			if node, err = gmc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(gmc.hooks) - 1; i >= 0; i-- {
			if gmc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = gmc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, gmc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (gmc *GroupMemberCreate) SaveX(ctx context.Context) *GroupMember {
	v, err := gmc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gmc *GroupMemberCreate) Exec(ctx context.Context) error {
	_, err := gmc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gmc *GroupMemberCreate) ExecX(ctx context.Context) {
	if err := gmc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (gmc *GroupMemberCreate) check() error {
	if _, ok := gmc.mutation.UID(); !ok {
		return &ValidationError{Name: "uid", err: errors.New(`ent: missing required field "GroupMember.uid"`)}
	}
	if _, ok := gmc.mutation.GroupID(); !ok {
		return &ValidationError{Name: "group_id", err: errors.New(`ent: missing required field "GroupMember.group_id"`)}
	}
	if _, ok := gmc.mutation.Inviter(); !ok {
		return &ValidationError{Name: "inviter", err: errors.New(`ent: missing required field "GroupMember.inviter"`)}
	}
	if _, ok := gmc.mutation.Remark(); !ok {
		return &ValidationError{Name: "remark", err: errors.New(`ent: missing required field "GroupMember.remark"`)}
	}
	if _, ok := gmc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "GroupMember.status"`)}
	}
	return nil
}

func (gmc *GroupMemberCreate) sqlSave(ctx context.Context) (*GroupMember, error) {
	_node, _spec := gmc.createSpec()
	if err := sqlgraph.CreateNode(ctx, gmc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	return _node, nil
}

func (gmc *GroupMemberCreate) createSpec() (*GroupMember, *sqlgraph.CreateSpec) {
	var (
		_node = &GroupMember{config: gmc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: groupmember.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: groupmember.FieldID,
			},
		}
	)
	if id, ok := gmc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := gmc.mutation.UID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: groupmember.FieldUID,
		})
		_node.UID = value
	}
	if value, ok := gmc.mutation.GroupID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: groupmember.FieldGroupID,
		})
		_node.GroupID = value
	}
	if value, ok := gmc.mutation.Inviter(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: groupmember.FieldInviter,
		})
		_node.Inviter = value
	}
	if value, ok := gmc.mutation.Remark(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: groupmember.FieldRemark,
		})
		_node.Remark = value
	}
	if value, ok := gmc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: groupmember.FieldStatus,
		})
		_node.Status = value
	}
	if value, ok := gmc.mutation.ApplyAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: groupmember.FieldApplyAt,
		})
		_node.ApplyAt = value
	}
	if value, ok := gmc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: groupmember.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := gmc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: groupmember.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := gmc.mutation.DeletedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: groupmember.FieldDeletedAt,
		})
		_node.DeletedAt = &value
	}
	return _node, _spec
}

// GroupMemberCreateBulk is the builder for creating many GroupMember entities in bulk.
type GroupMemberCreateBulk struct {
	config
	builders []*GroupMemberCreate
}

// Save creates the GroupMember entities in the database.
func (gmcb *GroupMemberCreateBulk) Save(ctx context.Context) ([]*GroupMember, error) {
	specs := make([]*sqlgraph.CreateSpec, len(gmcb.builders))
	nodes := make([]*GroupMember, len(gmcb.builders))
	mutators := make([]Mutator, len(gmcb.builders))
	for i := range gmcb.builders {
		func(i int, root context.Context) {
			builder := gmcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*GroupMemberMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, gmcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, gmcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int64(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, gmcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (gmcb *GroupMemberCreateBulk) SaveX(ctx context.Context) []*GroupMember {
	v, err := gmcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gmcb *GroupMemberCreateBulk) Exec(ctx context.Context) error {
	_, err := gmcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gmcb *GroupMemberCreateBulk) ExecX(ctx context.Context) {
	if err := gmcb.Exec(ctx); err != nil {
		panic(err)
	}
}
