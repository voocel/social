package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// GroupMember holds the schema definition for the GroupMember entity.
type GroupMember struct {
	ent.Schema
}

// Fields of the GroupMember.
func (GroupMember) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("uid"),
		field.Int64("group_id"),
		field.Int64("inviter").Default(0),
		field.String("remark").Default(""),
		field.Int8("status").Default(0),
		field.Time("apply_at").Optional().StructTag(`json:"-"`),
		field.Time("created_at").Optional().StructTag(`json:"-"`),
		field.Time("updated_at").Optional().StructTag(`json:"-"`),
		field.Time("deleted_at").Optional().Nillable().StructTag(`json:"-"`),
	}
}

// Edges of the GroupMember.
func (GroupMember) Edges() []ent.Edge {
	return nil
}
