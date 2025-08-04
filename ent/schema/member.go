package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"fmt"
	"github.com/geowa4/rara-membership/validation"
	"slices"
)

// Member holds the schema definition for the Member entity.
type Member struct {
	ent.Schema
}

// Fields of the Member.
func (Member) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("email").Validate(validation.ValidateEmail),
		field.String("phone").Optional(),
		field.String("mailing_address").Optional(),
		field.String("call_sign").Optional(),
		field.String("frn").Optional(),
		field.Bool("is_active").Default(true),
		field.Bool("is_silent_key").Default(false),
		field.String("license_class").
			Optional().
			Validate(func(c string) error {
				classes := []string{"Technician", "General", "Extra", "Novice", "Advanced"}
				if !slices.Contains(classes, c) {
					return fmt.Errorf("invalid license class %s is not one of %+q", c, classes)
				}
				return nil
			}),
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
