package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// PointDeduction holds the schema definition for the PointDeduction entity.
type PointDeduction struct {
	ent.Schema
}

// Fields of the PointDeduction.
func (PointDeduction) Fields() []ent.Field {
	return []ent.Field{
		field.Int("points").
			Min(1).
			Comment("Number of points being deducted from the member"),
		field.String("notes").
			Optional().
			Comment("Optional notes about the redemption"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("When the deduction was created"),
	}
}

// Edges of the PointDeduction.
func (PointDeduction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("member", Member.Type).
			Unique().
			Required().
			Comment("The member redeeming the points"),
	}
}
