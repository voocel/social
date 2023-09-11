package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Friend holds the schema definition for the Friend entity.
type Friend struct {
	ent.Schema
}

// Fields of the Friend.
func (Friend) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("uid"),
		field.Int64("friend_id"),
		field.String("remark"),
		field.Int8("shield").Default(0),
		field.Time("created_at").Optional().StructTag(`json:"-"`),
		field.Time("updated_at").Optional().StructTag(`json:"-"`),
		field.Time("deleted_at").Optional().Nillable().StructTag(`json:"-"`),
	}
}

// Edges of the Friend.
func (Friend) Edges() []ent.Edge {
	return nil
}
