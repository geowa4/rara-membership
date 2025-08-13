package schema

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	gen "github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/ent/hook"
	"github.com/geowa4/rara-membership/validation"
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
		field.String("member_type").
			Optional().
			Validate(func(t string) error {
				types := []string{"Student", "Regular", "Senior", "Associate"}
				if !slices.Contains(types, t) {
					return fmt.Errorf("invalid member type %s is not one of %+q", t, types)
				}
				return nil
			}),
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
		index.Fields("call_sign").
			Unique().
			Annotations(entsql.IndexWhere("is_active = true")),
	}
}

// Edges of the Member.
func (Member) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("point_allocations", PointAllocation.Type).
			Ref("member").
			Comment("Point allocations received by this member"),
		edge.From("point_deductions", PointDeduction.Type).
			Ref("member").
			Comment("Point deductions for redemptions by this member"),
	}
}

// Hooks of the Member.
func (Member) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(func(next ent.Mutator) ent.Mutator {
			return hook.MemberFunc(func(ctx context.Context, m *gen.MemberMutation) (ent.Value, error) {
				if callSign, ok := m.CallSign(); ok {
					m.SetCallSign(strings.ToUpper(callSign))
				}
				return next.Mutate(ctx, m)
			})
		}, ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne),
	}
	//return []ent.Hook{
	//	func(next ent.Mutator) ent.Mutator {
	//		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	//			// Only process Member mutations
	//			if memberMutation, ok := m.(interface {
	//				CallSign() (string, bool)
	//				SetCallSign(string)
	//			}); ok {
	//				// Convert call sign to uppercase if it's being set
	//				if callSign, exists := memberMutation.CallSign(); exists && callSign != "" {
	//					memberMutation.SetCallSign(strings.ToUpper(callSign))
	//				}
	//			}
	//			return next.Mutate(ctx, m)
	//		})
	//	},
	//}
}
