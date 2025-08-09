package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// PointAllocation holds the schema definition for the PointAllocation entity.
type PointAllocation struct {
	ent.Schema
}

// Fields of the PointAllocation.
func (PointAllocation) Fields() []ent.Field {
	return []ent.Field{
		field.Int("points").
			Min(0).
			Max(100).
			Comment("Number of points allocated to the member for the event"),
		field.String("notes").
			Optional().
			Comment("Optional notes about the allocation"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("When the allocation was created"),
	}
}

// Edges of the PointAllocation.
func (PointAllocation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("event", Event.Type).
			Unique().
			Required().
			Comment("The event for which points are allocated"),
		edge.To("member", Member.Type).
			Unique().
			Required().
			Comment("The member receiving the points"),
	}
}

// Indexes of the PointAllocation.
func (PointAllocation) Indexes() []ent.Index {
	return []ent.Index{
		// Ensure a member can only receive points once per event
		index.Edges("event", "member").
			Unique(),
	}
}
