package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Group holds the schema definition for the Group entity.
type Group struct {
	ent.Schema
}

// Fields of the Group.
func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("name"),
		field.Int64("owner"),
		field.String("avatar").Default(""),
		field.Int64("created_uid"),
		field.Int8("mode").Default(0),
		field.Int8("type").Default(0),
		field.Int8("status").Default(0),
		field.Int8("invite_mode").Default(0),
		field.String("notice").Default(""),
		field.String("introduction").Default(""),
		field.Time("created_at").Optional().StructTag(`json:"-"`),
		field.Time("updated_at").Optional().StructTag(`json:"-"`),
		field.Time("deleted_at").Optional().Nillable().StructTag(`json:"-"`),
	}
}

// Edges of the Friend.
func (Group) Edges() []ent.Edge {
	return nil
}
