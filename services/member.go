package services

import (
	"context"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/ent/member"
)

func GetActiveMembers() ([]*ent.Member, error) {
	ctx := context.Background()
	return database.Client.Member.Query().Where(member.IsActiveEQ(true)).All(ctx)
}

func CreateMember(name, email, phone, callSign, frn string, isActive bool) (*ent.Member, error) {
	ctx := context.Background()
	return database.Client.Member.Create().
		SetName(name).
		SetEmail(email).
		SetPhone(phone).
		SetCallSign(callSign).
		SetFrn(frn).
		SetIsActive(isActive).
		Save(ctx)
}
