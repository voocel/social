// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"social/ent/friend"
	"social/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// FriendUpdate is the builder for updating Friend entities.
type FriendUpdate struct {
	config
	hooks    []Hook
	mutation *FriendMutation
}

// Where appends a list predicates to the FriendUpdate builder.
func (fu *FriendUpdate) Where(ps ...predicate.Friend) *FriendUpdate {
	fu.mutation.Where(ps...)
	return fu
}

// SetUID sets the "uid" field.
func (fu *FriendUpdate) SetUID(i int64) *FriendUpdate {
	fu.mutation.ResetUID()
	fu.mutation.SetUID(i)
	return fu
}

// AddUID adds i to the "uid" field.
func (fu *FriendUpdate) AddUID(i int64) *FriendUpdate {
	fu.mutation.AddUID(i)
	return fu
}

// SetFriendID sets the "friend_id" field.
func (fu *FriendUpdate) SetFriendID(i int64) *FriendUpdate {
	fu.mutation.ResetFriendID()
	fu.mutation.SetFriendID(i)
	return fu
}

// AddFriendID adds i to the "friend_id" field.
func (fu *FriendUpdate) AddFriendID(i int64) *FriendUpdate {
	fu.mutation.AddFriendID(i)
	return fu
}

// SetRemark sets the "remark" field.
func (fu *FriendUpdate) SetRemark(s string) *FriendUpdate {
	fu.mutation.SetRemark(s)
	return fu
}

// SetShield sets the "shield" field.
func (fu *FriendUpdate) SetShield(i int8) *FriendUpdate {
	fu.mutation.ResetShield()
	fu.mutation.SetShield(i)
	return fu
}

// SetNillableShield sets the "shield" field if the given value is not nil.
func (fu *FriendUpdate) SetNillableShield(i *int8) *FriendUpdate {
	if i != nil {
		fu.SetShield(*i)
	}
	return fu
}

// AddShield adds i to the "shield" field.
func (fu *FriendUpdate) AddShield(i int8) *FriendUpdate {
	fu.mutation.AddShield(i)
	return fu
}

// SetCreatedAt sets the "created_at" field.
func (fu *FriendUpdate) SetCreatedAt(t time.Time) *FriendUpdate {
	fu.mutation.SetCreatedAt(t)
	return fu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (fu *FriendUpdate) SetNillableCreatedAt(t *time.Time) *FriendUpdate {
	if t != nil {
		fu.SetCreatedAt(*t)
	}
	return fu
}

// ClearCreatedAt clears the value of the "created_at" field.
func (fu *FriendUpdate) ClearCreatedAt() *FriendUpdate {
	fu.mutation.ClearCreatedAt()
	return fu
}

// SetUpdatedAt sets the "updated_at" field.
func (fu *FriendUpdate) SetUpdatedAt(t time.Time) *FriendUpdate {
	fu.mutation.SetUpdatedAt(t)
	return fu
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (fu *FriendUpdate) SetNillableUpdatedAt(t *time.Time) *FriendUpdate {
	if t != nil {
		fu.SetUpdatedAt(*t)
	}
	return fu
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (fu *FriendUpdate) ClearUpdatedAt() *FriendUpdate {
	fu.mutation.ClearUpdatedAt()
	return fu
}

// SetDeletedAt sets the "deleted_at" field.
func (fu *FriendUpdate) SetDeletedAt(t time.Time) *FriendUpdate {
	fu.mutation.SetDeletedAt(t)
	return fu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (fu *FriendUpdate) SetNillableDeletedAt(t *time.Time) *FriendUpdate {
	if t != nil {
		fu.SetDeletedAt(*t)
	}
	return fu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (fu *FriendUpdate) ClearDeletedAt() *FriendUpdate {
	fu.mutation.ClearDeletedAt()
	return fu
}

// Mutation returns the FriendMutation object of the builder.
func (fu *FriendUpdate) Mutation() *FriendMutation {
	return fu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fu *FriendUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(fu.hooks) == 0 {
		affected, err = fu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FriendMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fu.mutation = mutation
			affected, err = fu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(fu.hooks) - 1; i >= 0; i-- {
			if fu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = fu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (fu *FriendUpdate) SaveX(ctx context.Context) int {
	affected, err := fu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fu *FriendUpdate) Exec(ctx context.Context) error {
	_, err := fu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fu *FriendUpdate) ExecX(ctx context.Context) {
	if err := fu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fu *FriendUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   friend.Table,
			Columns: friend.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: friend.FieldID,
			},
		},
	}
	if ps := fu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fu.mutation.UID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friend.FieldUID,
		})
	}
	if value, ok := fu.mutation.AddedUID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friend.FieldUID,
		})
	}
	if value, ok := fu.mutation.FriendID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friend.FieldFriendID,
		})
	}
	if value, ok := fu.mutation.AddedFriendID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friend.FieldFriendID,
		})
	}
	if value, ok := fu.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: friend.FieldRemark,
		})
	}
	if value, ok := fu.mutation.Shield(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: friend.FieldShield,
		})
	}
	if value, ok := fu.mutation.AddedShield(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: friend.FieldShield,
		})
	}
	if value, ok := fu.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friend.FieldCreatedAt,
		})
	}
	if fu.mutation.CreatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friend.FieldCreatedAt,
		})
	}
	if value, ok := fu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friend.FieldUpdatedAt,
		})
	}
	if fu.mutation.UpdatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friend.FieldUpdatedAt,
		})
	}
	if value, ok := fu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friend.FieldDeletedAt,
		})
	}
	if fu.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friend.FieldDeletedAt,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{friend.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// FriendUpdateOne is the builder for updating a single Friend entity.
type FriendUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FriendMutation
}

// SetUID sets the "uid" field.
func (fuo *FriendUpdateOne) SetUID(i int64) *FriendUpdateOne {
	fuo.mutation.ResetUID()
	fuo.mutation.SetUID(i)
	return fuo
}

// AddUID adds i to the "uid" field.
func (fuo *FriendUpdateOne) AddUID(i int64) *FriendUpdateOne {
	fuo.mutation.AddUID(i)
	return fuo
}

// SetFriendID sets the "friend_id" field.
func (fuo *FriendUpdateOne) SetFriendID(i int64) *FriendUpdateOne {
	fuo.mutation.ResetFriendID()
	fuo.mutation.SetFriendID(i)
	return fuo
}

// AddFriendID adds i to the "friend_id" field.
func (fuo *FriendUpdateOne) AddFriendID(i int64) *FriendUpdateOne {
	fuo.mutation.AddFriendID(i)
	return fuo
}

// SetRemark sets the "remark" field.
func (fuo *FriendUpdateOne) SetRemark(s string) *FriendUpdateOne {
	fuo.mutation.SetRemark(s)
	return fuo
}

// SetShield sets the "shield" field.
func (fuo *FriendUpdateOne) SetShield(i int8) *FriendUpdateOne {
	fuo.mutation.ResetShield()
	fuo.mutation.SetShield(i)
	return fuo
}

// SetNillableShield sets the "shield" field if the given value is not nil.
func (fuo *FriendUpdateOne) SetNillableShield(i *int8) *FriendUpdateOne {
	if i != nil {
		fuo.SetShield(*i)
	}
	return fuo
}

// AddShield adds i to the "shield" field.
func (fuo *FriendUpdateOne) AddShield(i int8) *FriendUpdateOne {
	fuo.mutation.AddShield(i)
	return fuo
}

// SetCreatedAt sets the "created_at" field.
func (fuo *FriendUpdateOne) SetCreatedAt(t time.Time) *FriendUpdateOne {
	fuo.mutation.SetCreatedAt(t)
	return fuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (fuo *FriendUpdateOne) SetNillableCreatedAt(t *time.Time) *FriendUpdateOne {
	if t != nil {
		fuo.SetCreatedAt(*t)
	}
	return fuo
}

// ClearCreatedAt clears the value of the "created_at" field.
func (fuo *FriendUpdateOne) ClearCreatedAt() *FriendUpdateOne {
	fuo.mutation.ClearCreatedAt()
	return fuo
}

// SetUpdatedAt sets the "updated_at" field.
func (fuo *FriendUpdateOne) SetUpdatedAt(t time.Time) *FriendUpdateOne {
	fuo.mutation.SetUpdatedAt(t)
	return fuo
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (fuo *FriendUpdateOne) SetNillableUpdatedAt(t *time.Time) *FriendUpdateOne {
	if t != nil {
		fuo.SetUpdatedAt(*t)
	}
	return fuo
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (fuo *FriendUpdateOne) ClearUpdatedAt() *FriendUpdateOne {
	fuo.mutation.ClearUpdatedAt()
	return fuo
}

// SetDeletedAt sets the "deleted_at" field.
func (fuo *FriendUpdateOne) SetDeletedAt(t time.Time) *FriendUpdateOne {
	fuo.mutation.SetDeletedAt(t)
	return fuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (fuo *FriendUpdateOne) SetNillableDeletedAt(t *time.Time) *FriendUpdateOne {
	if t != nil {
		fuo.SetDeletedAt(*t)
	}
	return fuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (fuo *FriendUpdateOne) ClearDeletedAt() *FriendUpdateOne {
	fuo.mutation.ClearDeletedAt()
	return fuo
}

// Mutation returns the FriendMutation object of the builder.
func (fuo *FriendUpdateOne) Mutation() *FriendMutation {
	return fuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fuo *FriendUpdateOne) Select(field string, fields ...string) *FriendUpdateOne {
	fuo.fields = append([]string{field}, fields...)
	return fuo
}

// Save executes the query and returns the updated Friend entity.
func (fuo *FriendUpdateOne) Save(ctx context.Context) (*Friend, error) {
	var (
		err  error
		node *Friend
	)
	if len(fuo.hooks) == 0 {
		node, err = fuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FriendMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fuo.mutation = mutation
			node, err = fuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(fuo.hooks) - 1; i >= 0; i-- {
			if fuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = fuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (fuo *FriendUpdateOne) SaveX(ctx context.Context) *Friend {
	node, err := fuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fuo *FriendUpdateOne) Exec(ctx context.Context) error {
	_, err := fuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuo *FriendUpdateOne) ExecX(ctx context.Context) {
	if err := fuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fuo *FriendUpdateOne) sqlSave(ctx context.Context) (_node *Friend, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   friend.Table,
			Columns: friend.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: friend.FieldID,
			},
		},
	}
	id, ok := fuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Friend.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, friend.FieldID)
		for _, f := range fields {
			if !friend.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != friend.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fuo.mutation.UID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friend.FieldUID,
		})
	}
	if value, ok := fuo.mutation.AddedUID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friend.FieldUID,
		})
	}
	if value, ok := fuo.mutation.FriendID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friend.FieldFriendID,
		})
	}
	if value, ok := fuo.mutation.AddedFriendID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friend.FieldFriendID,
		})
	}
	if value, ok := fuo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: friend.FieldRemark,
		})
	}
	if value, ok := fuo.mutation.Shield(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: friend.FieldShield,
		})
	}
	if value, ok := fuo.mutation.AddedShield(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: friend.FieldShield,
		})
	}
	if value, ok := fuo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friend.FieldCreatedAt,
		})
	}
	if fuo.mutation.CreatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friend.FieldCreatedAt,
		})
	}
	if value, ok := fuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friend.FieldUpdatedAt,
		})
	}
	if fuo.mutation.UpdatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friend.FieldUpdatedAt,
		})
	}
	if value, ok := fuo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friend.FieldDeletedAt,
		})
	}
	if fuo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friend.FieldDeletedAt,
		})
	}
	_node = &Friend{config: fuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{friend.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
