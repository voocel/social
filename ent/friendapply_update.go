// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"social/ent/friendapply"
	"social/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// FriendApplyUpdate is the builder for updating FriendApply entities.
type FriendApplyUpdate struct {
	config
	hooks    []Hook
	mutation *FriendApplyMutation
}

// Where appends a list predicates to the FriendApplyUpdate builder.
func (fau *FriendApplyUpdate) Where(ps ...predicate.FriendApply) *FriendApplyUpdate {
	fau.mutation.Where(ps...)
	return fau
}

// SetFromID sets the "from_id" field.
func (fau *FriendApplyUpdate) SetFromID(i int64) *FriendApplyUpdate {
	fau.mutation.ResetFromID()
	fau.mutation.SetFromID(i)
	return fau
}

// AddFromID adds i to the "from_id" field.
func (fau *FriendApplyUpdate) AddFromID(i int64) *FriendApplyUpdate {
	fau.mutation.AddFromID(i)
	return fau
}

// SetToID sets the "to_id" field.
func (fau *FriendApplyUpdate) SetToID(i int64) *FriendApplyUpdate {
	fau.mutation.ResetToID()
	fau.mutation.SetToID(i)
	return fau
}

// AddToID adds i to the "to_id" field.
func (fau *FriendApplyUpdate) AddToID(i int64) *FriendApplyUpdate {
	fau.mutation.AddToID(i)
	return fau
}

// SetRemark sets the "remark" field.
func (fau *FriendApplyUpdate) SetRemark(s string) *FriendApplyUpdate {
	fau.mutation.SetRemark(s)
	return fau
}

// SetStatus sets the "status" field.
func (fau *FriendApplyUpdate) SetStatus(i int8) *FriendApplyUpdate {
	fau.mutation.ResetStatus()
	fau.mutation.SetStatus(i)
	return fau
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (fau *FriendApplyUpdate) SetNillableStatus(i *int8) *FriendApplyUpdate {
	if i != nil {
		fau.SetStatus(*i)
	}
	return fau
}

// AddStatus adds i to the "status" field.
func (fau *FriendApplyUpdate) AddStatus(i int8) *FriendApplyUpdate {
	fau.mutation.AddStatus(i)
	return fau
}

// ClearStatus clears the value of the "status" field.
func (fau *FriendApplyUpdate) ClearStatus() *FriendApplyUpdate {
	fau.mutation.ClearStatus()
	return fau
}

// SetCreatedAt sets the "created_at" field.
func (fau *FriendApplyUpdate) SetCreatedAt(t time.Time) *FriendApplyUpdate {
	fau.mutation.SetCreatedAt(t)
	return fau
}

// SetUpdatedAt sets the "updated_at" field.
func (fau *FriendApplyUpdate) SetUpdatedAt(t time.Time) *FriendApplyUpdate {
	fau.mutation.SetUpdatedAt(t)
	return fau
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (fau *FriendApplyUpdate) SetNillableUpdatedAt(t *time.Time) *FriendApplyUpdate {
	if t != nil {
		fau.SetUpdatedAt(*t)
	}
	return fau
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (fau *FriendApplyUpdate) ClearUpdatedAt() *FriendApplyUpdate {
	fau.mutation.ClearUpdatedAt()
	return fau
}

// SetDeletedAt sets the "deleted_at" field.
func (fau *FriendApplyUpdate) SetDeletedAt(t time.Time) *FriendApplyUpdate {
	fau.mutation.SetDeletedAt(t)
	return fau
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (fau *FriendApplyUpdate) SetNillableDeletedAt(t *time.Time) *FriendApplyUpdate {
	if t != nil {
		fau.SetDeletedAt(*t)
	}
	return fau
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (fau *FriendApplyUpdate) ClearDeletedAt() *FriendApplyUpdate {
	fau.mutation.ClearDeletedAt()
	return fau
}

// Mutation returns the FriendApplyMutation object of the builder.
func (fau *FriendApplyUpdate) Mutation() *FriendApplyMutation {
	return fau.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fau *FriendApplyUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(fau.hooks) == 0 {
		affected, err = fau.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FriendApplyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fau.mutation = mutation
			affected, err = fau.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(fau.hooks) - 1; i >= 0; i-- {
			if fau.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = fau.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fau.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (fau *FriendApplyUpdate) SaveX(ctx context.Context) int {
	affected, err := fau.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fau *FriendApplyUpdate) Exec(ctx context.Context) error {
	_, err := fau.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fau *FriendApplyUpdate) ExecX(ctx context.Context) {
	if err := fau.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fau *FriendApplyUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   friendapply.Table,
			Columns: friendapply.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: friendapply.FieldID,
			},
		},
	}
	if ps := fau.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fau.mutation.FromID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friendapply.FieldFromID,
		})
	}
	if value, ok := fau.mutation.AddedFromID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friendapply.FieldFromID,
		})
	}
	if value, ok := fau.mutation.ToID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friendapply.FieldToID,
		})
	}
	if value, ok := fau.mutation.AddedToID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friendapply.FieldToID,
		})
	}
	if value, ok := fau.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: friendapply.FieldRemark,
		})
	}
	if value, ok := fau.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: friendapply.FieldStatus,
		})
	}
	if value, ok := fau.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: friendapply.FieldStatus,
		})
	}
	if fau.mutation.StatusCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Column: friendapply.FieldStatus,
		})
	}
	if value, ok := fau.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friendapply.FieldCreatedAt,
		})
	}
	if value, ok := fau.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friendapply.FieldUpdatedAt,
		})
	}
	if fau.mutation.UpdatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friendapply.FieldUpdatedAt,
		})
	}
	if value, ok := fau.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friendapply.FieldDeletedAt,
		})
	}
	if fau.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friendapply.FieldDeletedAt,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fau.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{friendapply.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// FriendApplyUpdateOne is the builder for updating a single FriendApply entity.
type FriendApplyUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FriendApplyMutation
}

// SetFromID sets the "from_id" field.
func (fauo *FriendApplyUpdateOne) SetFromID(i int64) *FriendApplyUpdateOne {
	fauo.mutation.ResetFromID()
	fauo.mutation.SetFromID(i)
	return fauo
}

// AddFromID adds i to the "from_id" field.
func (fauo *FriendApplyUpdateOne) AddFromID(i int64) *FriendApplyUpdateOne {
	fauo.mutation.AddFromID(i)
	return fauo
}

// SetToID sets the "to_id" field.
func (fauo *FriendApplyUpdateOne) SetToID(i int64) *FriendApplyUpdateOne {
	fauo.mutation.ResetToID()
	fauo.mutation.SetToID(i)
	return fauo
}

// AddToID adds i to the "to_id" field.
func (fauo *FriendApplyUpdateOne) AddToID(i int64) *FriendApplyUpdateOne {
	fauo.mutation.AddToID(i)
	return fauo
}

// SetRemark sets the "remark" field.
func (fauo *FriendApplyUpdateOne) SetRemark(s string) *FriendApplyUpdateOne {
	fauo.mutation.SetRemark(s)
	return fauo
}

// SetStatus sets the "status" field.
func (fauo *FriendApplyUpdateOne) SetStatus(i int8) *FriendApplyUpdateOne {
	fauo.mutation.ResetStatus()
	fauo.mutation.SetStatus(i)
	return fauo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (fauo *FriendApplyUpdateOne) SetNillableStatus(i *int8) *FriendApplyUpdateOne {
	if i != nil {
		fauo.SetStatus(*i)
	}
	return fauo
}

// AddStatus adds i to the "status" field.
func (fauo *FriendApplyUpdateOne) AddStatus(i int8) *FriendApplyUpdateOne {
	fauo.mutation.AddStatus(i)
	return fauo
}

// ClearStatus clears the value of the "status" field.
func (fauo *FriendApplyUpdateOne) ClearStatus() *FriendApplyUpdateOne {
	fauo.mutation.ClearStatus()
	return fauo
}

// SetCreatedAt sets the "created_at" field.
func (fauo *FriendApplyUpdateOne) SetCreatedAt(t time.Time) *FriendApplyUpdateOne {
	fauo.mutation.SetCreatedAt(t)
	return fauo
}

// SetUpdatedAt sets the "updated_at" field.
func (fauo *FriendApplyUpdateOne) SetUpdatedAt(t time.Time) *FriendApplyUpdateOne {
	fauo.mutation.SetUpdatedAt(t)
	return fauo
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (fauo *FriendApplyUpdateOne) SetNillableUpdatedAt(t *time.Time) *FriendApplyUpdateOne {
	if t != nil {
		fauo.SetUpdatedAt(*t)
	}
	return fauo
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (fauo *FriendApplyUpdateOne) ClearUpdatedAt() *FriendApplyUpdateOne {
	fauo.mutation.ClearUpdatedAt()
	return fauo
}

// SetDeletedAt sets the "deleted_at" field.
func (fauo *FriendApplyUpdateOne) SetDeletedAt(t time.Time) *FriendApplyUpdateOne {
	fauo.mutation.SetDeletedAt(t)
	return fauo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (fauo *FriendApplyUpdateOne) SetNillableDeletedAt(t *time.Time) *FriendApplyUpdateOne {
	if t != nil {
		fauo.SetDeletedAt(*t)
	}
	return fauo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (fauo *FriendApplyUpdateOne) ClearDeletedAt() *FriendApplyUpdateOne {
	fauo.mutation.ClearDeletedAt()
	return fauo
}

// Mutation returns the FriendApplyMutation object of the builder.
func (fauo *FriendApplyUpdateOne) Mutation() *FriendApplyMutation {
	return fauo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fauo *FriendApplyUpdateOne) Select(field string, fields ...string) *FriendApplyUpdateOne {
	fauo.fields = append([]string{field}, fields...)
	return fauo
}

// Save executes the query and returns the updated FriendApply entity.
func (fauo *FriendApplyUpdateOne) Save(ctx context.Context) (*FriendApply, error) {
	var (
		err  error
		node *FriendApply
	)
	if len(fauo.hooks) == 0 {
		node, err = fauo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FriendApplyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fauo.mutation = mutation
			node, err = fauo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(fauo.hooks) - 1; i >= 0; i-- {
			if fauo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = fauo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fauo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (fauo *FriendApplyUpdateOne) SaveX(ctx context.Context) *FriendApply {
	node, err := fauo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fauo *FriendApplyUpdateOne) Exec(ctx context.Context) error {
	_, err := fauo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fauo *FriendApplyUpdateOne) ExecX(ctx context.Context) {
	if err := fauo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fauo *FriendApplyUpdateOne) sqlSave(ctx context.Context) (_node *FriendApply, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   friendapply.Table,
			Columns: friendapply.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: friendapply.FieldID,
			},
		},
	}
	id, ok := fauo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "FriendApply.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fauo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, friendapply.FieldID)
		for _, f := range fields {
			if !friendapply.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != friendapply.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fauo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fauo.mutation.FromID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friendapply.FieldFromID,
		})
	}
	if value, ok := fauo.mutation.AddedFromID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friendapply.FieldFromID,
		})
	}
	if value, ok := fauo.mutation.ToID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friendapply.FieldToID,
		})
	}
	if value, ok := fauo.mutation.AddedToID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: friendapply.FieldToID,
		})
	}
	if value, ok := fauo.mutation.Remark(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: friendapply.FieldRemark,
		})
	}
	if value, ok := fauo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: friendapply.FieldStatus,
		})
	}
	if value, ok := fauo.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: friendapply.FieldStatus,
		})
	}
	if fauo.mutation.StatusCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Column: friendapply.FieldStatus,
		})
	}
	if value, ok := fauo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friendapply.FieldCreatedAt,
		})
	}
	if value, ok := fauo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friendapply.FieldUpdatedAt,
		})
	}
	if fauo.mutation.UpdatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friendapply.FieldUpdatedAt,
		})
	}
	if value, ok := fauo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: friendapply.FieldDeletedAt,
		})
	}
	if fauo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: friendapply.FieldDeletedAt,
		})
	}
	_node = &FriendApply{config: fauo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fauo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{friendapply.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
