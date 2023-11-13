package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("username"),
		field.String("password"),
		field.String("mobile"),
		field.String("nickname"),
		field.String("email"),
		field.String("ip"),
		field.String("address"),
		field.String("avatar").Optional(),
		field.String("summary").Optional(),
		field.Int8("sex").Optional(),
		field.Int8("status").Optional(),
		field.Time("birthday").Optional().StructTag(`json:"-"`),
		field.Time("created_at").Optional().StructTag(`json:"-"`),
		field.Time("updated_at").Optional().StructTag(`json:"-"`),
		field.Time("deleted_at").Optional().Nillable().StructTag(`json:"-"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
