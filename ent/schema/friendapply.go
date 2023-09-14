package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// FriendApply holds the schema definition for the FriendApply entity.
type FriendApply struct {
	ent.Schema
}

// Fields of the FriendApply.
func (FriendApply) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("from_id"),
		field.Int64("to_id"),
		field.String("remark"),
		field.Int8("status").Optional().StructTag(`json:"status"`),
		field.Time("created_at"),
		field.Time("updated_at").Optional().StructTag(`json:"-"`),
		field.Time("deleted_at").Optional().Nillable().StructTag(`json:"-"`),
	}
}

// Edges of the FriendApply.
func (FriendApply) Edges() []ent.Edge {
	return nil
}
