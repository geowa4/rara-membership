package services

import (
	"context"
	"strings"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/ent/member"
)

func GetActiveMembers() ([]*ent.Member, error) {
	ctx := context.Background()
	return database.Client.Member.Query().Where(member.IsActiveEQ(true)).All(ctx)
}

func CreateMember(name, email, phone, mailingAddress, callSign, frn string, isActive, isSilentKey bool, licenseClass string) (*ent.Member, error) {
	ctx := context.Background()
	create := database.Client.Member.Create().
		SetName(name).
		SetEmail(email).
		SetIsActive(isActive).
		SetIsSilentKey(isSilentKey)

	if phone != "" {
		create = create.SetPhone(phone)
	}

	if mailingAddress != "" {
		create = create.SetMailingAddress(mailingAddress)
	}

	if callSign != "" {
		create = create.SetCallSign(callSign)
	}

	if frn != "" {
		create = create.SetFrn(frn)
	}

	if licenseClass != "" {
		create = create.SetLicenseClass(licenseClass)
	}

	return create.Save(ctx)
}

func GetActiveMemberByCallSign(callSign string) (*ent.Member, error) {
	ctx := context.Background()
	return database.Client.Member.Query().
		Where(
			member.Or(member.CallSignEQ(strings.ToUpper(callSign)), member.CallSignEQ(strings.ToLower(callSign))),
			member.IsActiveEQ(true),
		).
		Only(ctx)
}

func UpdateMember(id int, name, email, phone, mailingAddress, callSign, frn string, isActive, isSilentKey bool, licenseClass string) (*ent.Member, error) {
	ctx := context.Background()
	update := database.Client.Member.UpdateOneID(id).
		SetName(name).
		SetEmail(email).
		SetIsActive(isActive).
		SetIsSilentKey(isSilentKey)

	if phone != "" {
		update = update.SetPhone(phone)
	} else {
		update = update.ClearPhone()
	}

	if mailingAddress != "" {
		update = update.SetMailingAddress(mailingAddress)
	} else {
		update = update.ClearMailingAddress()
	}

	if callSign != "" {
		update = update.SetCallSign(callSign)
	} else {
		update = update.ClearCallSign()
	}

	if frn != "" {
		update = update.SetFrn(frn)
	} else {
		update = update.ClearFrn()
	}

	if licenseClass != "" {
		update = update.SetLicenseClass(licenseClass)
	} else {
		update = update.ClearLicenseClass()
	}

	return update.Save(ctx)
}
