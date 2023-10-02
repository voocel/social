// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"social/ent/message"
	"social/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MessageUpdate is the builder for updating Message entities.
type MessageUpdate struct {
	config
	hooks    []Hook
	mutation *MessageMutation
}

// Where appends a list predicates to the MessageUpdate builder.
func (mu *MessageUpdate) Where(ps ...predicate.Message) *MessageUpdate {
	mu.mutation.Where(ps...)
	return mu
}

// SetSenderID sets the "sender_id" field.
func (mu *MessageUpdate) SetSenderID(i int64) *MessageUpdate {
	mu.mutation.ResetSenderID()
	mu.mutation.SetSenderID(i)
	return mu
}

// AddSenderID adds i to the "sender_id" field.
func (mu *MessageUpdate) AddSenderID(i int64) *MessageUpdate {
	mu.mutation.AddSenderID(i)
	return mu
}

// SetReceiverID sets the "receiver_id" field.
func (mu *MessageUpdate) SetReceiverID(i int64) *MessageUpdate {
	mu.mutation.ResetReceiverID()
	mu.mutation.SetReceiverID(i)
	return mu
}

// AddReceiverID adds i to the "receiver_id" field.
func (mu *MessageUpdate) AddReceiverID(i int64) *MessageUpdate {
	mu.mutation.AddReceiverID(i)
	return mu
}

// SetContent sets the "content" field.
func (mu *MessageUpdate) SetContent(s string) *MessageUpdate {
	mu.mutation.SetContent(s)
	return mu
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableContent(s *string) *MessageUpdate {
	if s != nil {
		mu.SetContent(*s)
	}
	return mu
}

// SetStatus sets the "status" field.
func (mu *MessageUpdate) SetStatus(i int8) *MessageUpdate {
	mu.mutation.ResetStatus()
	mu.mutation.SetStatus(i)
	return mu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableStatus(i *int8) *MessageUpdate {
	if i != nil {
		mu.SetStatus(*i)
	}
	return mu
}

// AddStatus adds i to the "status" field.
func (mu *MessageUpdate) AddStatus(i int8) *MessageUpdate {
	mu.mutation.AddStatus(i)
	return mu
}

// SetCreatedAt sets the "created_at" field.
func (mu *MessageUpdate) SetCreatedAt(t time.Time) *MessageUpdate {
	mu.mutation.SetCreatedAt(t)
	return mu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableCreatedAt(t *time.Time) *MessageUpdate {
	if t != nil {
		mu.SetCreatedAt(*t)
	}
	return mu
}

// ClearCreatedAt clears the value of the "created_at" field.
func (mu *MessageUpdate) ClearCreatedAt() *MessageUpdate {
	mu.mutation.ClearCreatedAt()
	return mu
}

// SetUpdatedAt sets the "updated_at" field.
func (mu *MessageUpdate) SetUpdatedAt(t time.Time) *MessageUpdate {
	mu.mutation.SetUpdatedAt(t)
	return mu
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableUpdatedAt(t *time.Time) *MessageUpdate {
	if t != nil {
		mu.SetUpdatedAt(*t)
	}
	return mu
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (mu *MessageUpdate) ClearUpdatedAt() *MessageUpdate {
	mu.mutation.ClearUpdatedAt()
	return mu
}

// SetDeletedAt sets the "deleted_at" field.
func (mu *MessageUpdate) SetDeletedAt(t time.Time) *MessageUpdate {
	mu.mutation.SetDeletedAt(t)
	return mu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableDeletedAt(t *time.Time) *MessageUpdate {
	if t != nil {
		mu.SetDeletedAt(*t)
	}
	return mu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (mu *MessageUpdate) ClearDeletedAt() *MessageUpdate {
	mu.mutation.ClearDeletedAt()
	return mu
}

// Mutation returns the MessageMutation object of the builder.
func (mu *MessageUpdate) Mutation() *MessageMutation {
	return mu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mu *MessageUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(mu.hooks) == 0 {
		affected, err = mu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MessageMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			mu.mutation = mutation
			affected, err = mu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(mu.hooks) - 1; i >= 0; i-- {
			if mu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = mu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, mu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (mu *MessageUpdate) SaveX(ctx context.Context) int {
	affected, err := mu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mu *MessageUpdate) Exec(ctx context.Context) error {
	_, err := mu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mu *MessageUpdate) ExecX(ctx context.Context) {
	if err := mu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (mu *MessageUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   message.Table,
			Columns: message.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: message.FieldID,
			},
		},
	}
	if ps := mu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mu.mutation.SenderID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: message.FieldSenderID,
		})
	}
	if value, ok := mu.mutation.AddedSenderID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: message.FieldSenderID,
		})
	}
	if value, ok := mu.mutation.ReceiverID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: message.FieldReceiverID,
		})
	}
	if value, ok := mu.mutation.AddedReceiverID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: message.FieldReceiverID,
		})
	}
	if value, ok := mu.mutation.Content(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: message.FieldContent,
		})
	}
	if value, ok := mu.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: message.FieldStatus,
		})
	}
	if value, ok := mu.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: message.FieldStatus,
		})
	}
	if value, ok := mu.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: message.FieldCreatedAt,
		})
	}
	if mu.mutation.CreatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: message.FieldCreatedAt,
		})
	}
	if value, ok := mu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: message.FieldUpdatedAt,
		})
	}
	if mu.mutation.UpdatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: message.FieldUpdatedAt,
		})
	}
	if value, ok := mu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: message.FieldDeletedAt,
		})
	}
	if mu.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: message.FieldDeletedAt,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, mu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{message.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// MessageUpdateOne is the builder for updating a single Message entity.
type MessageUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *MessageMutation
}

// SetSenderID sets the "sender_id" field.
func (muo *MessageUpdateOne) SetSenderID(i int64) *MessageUpdateOne {
	muo.mutation.ResetSenderID()
	muo.mutation.SetSenderID(i)
	return muo
}

// AddSenderID adds i to the "sender_id" field.
func (muo *MessageUpdateOne) AddSenderID(i int64) *MessageUpdateOne {
	muo.mutation.AddSenderID(i)
	return muo
}

// SetReceiverID sets the "receiver_id" field.
func (muo *MessageUpdateOne) SetReceiverID(i int64) *MessageUpdateOne {
	muo.mutation.ResetReceiverID()
	muo.mutation.SetReceiverID(i)
	return muo
}

// AddReceiverID adds i to the "receiver_id" field.
func (muo *MessageUpdateOne) AddReceiverID(i int64) *MessageUpdateOne {
	muo.mutation.AddReceiverID(i)
	return muo
}

// SetContent sets the "content" field.
func (muo *MessageUpdateOne) SetContent(s string) *MessageUpdateOne {
	muo.mutation.SetContent(s)
	return muo
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableContent(s *string) *MessageUpdateOne {
	if s != nil {
		muo.SetContent(*s)
	}
	return muo
}

// SetStatus sets the "status" field.
func (muo *MessageUpdateOne) SetStatus(i int8) *MessageUpdateOne {
	muo.mutation.ResetStatus()
	muo.mutation.SetStatus(i)
	return muo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableStatus(i *int8) *MessageUpdateOne {
	if i != nil {
		muo.SetStatus(*i)
	}
	return muo
}

// AddStatus adds i to the "status" field.
func (muo *MessageUpdateOne) AddStatus(i int8) *MessageUpdateOne {
	muo.mutation.AddStatus(i)
	return muo
}

// SetCreatedAt sets the "created_at" field.
func (muo *MessageUpdateOne) SetCreatedAt(t time.Time) *MessageUpdateOne {
	muo.mutation.SetCreatedAt(t)
	return muo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableCreatedAt(t *time.Time) *MessageUpdateOne {
	if t != nil {
		muo.SetCreatedAt(*t)
	}
	return muo
}

// ClearCreatedAt clears the value of the "created_at" field.
func (muo *MessageUpdateOne) ClearCreatedAt() *MessageUpdateOne {
	muo.mutation.ClearCreatedAt()
	return muo
}

// SetUpdatedAt sets the "updated_at" field.
func (muo *MessageUpdateOne) SetUpdatedAt(t time.Time) *MessageUpdateOne {
	muo.mutation.SetUpdatedAt(t)
	return muo
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableUpdatedAt(t *time.Time) *MessageUpdateOne {
	if t != nil {
		muo.SetUpdatedAt(*t)
	}
	return muo
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (muo *MessageUpdateOne) ClearUpdatedAt() *MessageUpdateOne {
	muo.mutation.ClearUpdatedAt()
	return muo
}

// SetDeletedAt sets the "deleted_at" field.
func (muo *MessageUpdateOne) SetDeletedAt(t time.Time) *MessageUpdateOne {
	muo.mutation.SetDeletedAt(t)
	return muo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableDeletedAt(t *time.Time) *MessageUpdateOne {
	if t != nil {
		muo.SetDeletedAt(*t)
	}
	return muo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (muo *MessageUpdateOne) ClearDeletedAt() *MessageUpdateOne {
	muo.mutation.ClearDeletedAt()
	return muo
}

// Mutation returns the MessageMutation object of the builder.
func (muo *MessageUpdateOne) Mutation() *MessageMutation {
	return muo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (muo *MessageUpdateOne) Select(field string, fields ...string) *MessageUpdateOne {
	muo.fields = append([]string{field}, fields...)
	return muo
}

// Save executes the query and returns the updated Message entity.
func (muo *MessageUpdateOne) Save(ctx context.Context) (*Message, error) {
	var (
		err  error
		node *Message
	)
	if len(muo.hooks) == 0 {
		node, err = muo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MessageMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			muo.mutation = mutation
			node, err = muo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(muo.hooks) - 1; i >= 0; i-- {
			if muo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = muo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, muo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (muo *MessageUpdateOne) SaveX(ctx context.Context) *Message {
	node, err := muo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (muo *MessageUpdateOne) Exec(ctx context.Context) error {
	_, err := muo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (muo *MessageUpdateOne) ExecX(ctx context.Context) {
	if err := muo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (muo *MessageUpdateOne) sqlSave(ctx context.Context) (_node *Message, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   message.Table,
			Columns: message.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: message.FieldID,
			},
		},
	}
	id, ok := muo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Message.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := muo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, message.FieldID)
		for _, f := range fields {
			if !message.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != message.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := muo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := muo.mutation.SenderID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: message.FieldSenderID,
		})
	}
	if value, ok := muo.mutation.AddedSenderID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: message.FieldSenderID,
		})
	}
	if value, ok := muo.mutation.ReceiverID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: message.FieldReceiverID,
		})
	}
	if value, ok := muo.mutation.AddedReceiverID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: message.FieldReceiverID,
		})
	}
	if value, ok := muo.mutation.Content(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: message.FieldContent,
		})
	}
	if value, ok := muo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: message.FieldStatus,
		})
	}
	if value, ok := muo.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt8,
			Value:  value,
			Column: message.FieldStatus,
		})
	}
	if value, ok := muo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: message.FieldCreatedAt,
		})
	}
	if muo.mutation.CreatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: message.FieldCreatedAt,
		})
	}
	if value, ok := muo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: message.FieldUpdatedAt,
		})
	}
	if muo.mutation.UpdatedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: message.FieldUpdatedAt,
		})
	}
	if value, ok := muo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: message.FieldDeletedAt,
		})
	}
	if muo.mutation.DeletedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: message.FieldDeletedAt,
		})
	}
	_node = &Message{config: muo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, muo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{message.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
