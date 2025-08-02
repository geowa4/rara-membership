package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Member holds the schema definition for the Member entity.
type Member struct {
	ent.Schema
}

// Fields of the Member.
func (Member) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("email"),
		field.String("phone"),
		field.String("call_sign"),
		field.String("frn"),
		field.Bool("is_active").Default(true),
		field.Bool("is_silent_key").Default(false),
	}
}

// Indexes of the Member
func (Member) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("call_sign", "is_active"),
	}
}

// Edges of the Member.
func (Member) Edges() []ent.Edge {
	return nil
}
