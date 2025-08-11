package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

// Fields of the Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("description").NotEmpty(),
		field.Time("date").Default(time.Now),
		field.String("timezone").
			Default("America/New_York").
			Comment("IANA timezone identifier for the event"),
		field.Float("latitude").
			Default(43.138).
			Min(-90).
			Max(90),
		field.Float("longitude").
			Default(-77.572).
			Min(-180).
			Max(180),
		field.Int8("default_points_allocated").
			Default(10).
			Min(0).
			Max(100),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("point_allocations", PointAllocation.Type).
			Ref("event").
			Comment("Point allocations for this event"),
	}
}

// Indexes of the Event.
func (Event) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("date").
			Annotations(entsql.Desc()),
	}
}
