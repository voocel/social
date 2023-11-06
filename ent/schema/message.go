package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("sender_id"),
		field.Int64("receiver_id"),
		field.String("content").Default(""),
		field.Int8("content_type").Default(0),
		field.Int8("status").Default(0),
		field.Time("created_at").Optional().StructTag(`json:"-"`),
		field.Time("updated_at").Optional().StructTag(`json:"-"`),
		field.Time("deleted_at").Optional().Nillable().StructTag(`json:"-"`),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return nil
}
