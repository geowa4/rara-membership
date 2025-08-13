package services

import (
	"context"
	"strings"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/ent/member"
)

type MemberService struct {
	client *ent.Client
}

func NewMemberService(client *ent.Client) *MemberService {
	return &MemberService{client: client}
}

type CreateMemberInput struct {
	Name           string
	CallSign       string
	Email          string
	Phone          *string
	MailingAddress *string
	FRN            *string
	LicenseClass   *string
	MemberType     *string
}

type UpdateMemberInput struct {
	Name           *string
	CallSign       *string
	Email          *string
	Phone          *string
	MailingAddress *string
	FRN            *string
	LicenseClass   *string
	MemberType     *string
	IsActive       *bool
	IsSilentKey    *bool
}

func GetActiveMembers() ([]*ent.Member, error) {
	ctx := context.Background()
	return database.Client.Member.Query().Where(member.IsActiveEQ(true)).All(ctx)
}

func CreateMember(name, email, phone, mailingAddress, callSign, frn string, isActive, isSilentKey bool, licenseClass, memberType string) (*ent.Member, error) {
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

	if memberType != "" {
		create = create.SetMemberType(memberType)
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

func UpdateMember(id int, name, email, phone, mailingAddress, callSign, frn string, isActive, isSilentKey bool, licenseClass, memberType string) (*ent.Member, error) {
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

	if memberType != "" {
		update = update.SetMemberType(memberType)
	} else {
		update = update.ClearMemberType()
	}

	return update.Save(ctx)
}

func (s *MemberService) CreateMember(ctx context.Context, input CreateMemberInput) (*ent.Member, error) {
	create := s.client.Member.Create().
		SetName(input.Name).
		SetCallSign(strings.ToUpper(input.CallSign)).
		SetEmail(input.Email).
		SetIsActive(true).
		SetIsSilentKey(false)

	if input.Phone != nil && *input.Phone != "" {
		create = create.SetPhone(*input.Phone)
	}

	if input.MailingAddress != nil && *input.MailingAddress != "" {
		create = create.SetMailingAddress(*input.MailingAddress)
	}

	if input.FRN != nil && *input.FRN != "" {
		create = create.SetFrn(*input.FRN)
	}

	if input.LicenseClass != nil && *input.LicenseClass != "" {
		create = create.SetLicenseClass(*input.LicenseClass)
	}

	if input.MemberType != nil && *input.MemberType != "" {
		create = create.SetMemberType(*input.MemberType)
	}

	return create.Save(ctx)
}

func (s *MemberService) UpdateMemberByCallSign(ctx context.Context, callSign string, input UpdateMemberInput) (*ent.Member, error) {
	// First find the member
	mbr, err := s.client.Member.Query().
		Where(
			member.Or(
				member.CallSignEQ(strings.ToUpper(callSign)),
				member.CallSignEQ(strings.ToLower(callSign)),
			),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	update := s.client.Member.UpdateOneID(mbr.ID)

	if input.Name != nil {
		update = update.SetName(*input.Name)
	}

	if input.CallSign != nil {
		update = update.SetCallSign(strings.ToUpper(*input.CallSign))
	}

	if input.Email != nil {
		update = update.SetEmail(*input.Email)
	}

	if input.Phone != nil {
		if *input.Phone != "" {
			update = update.SetPhone(*input.Phone)
		} else {
			update = update.ClearPhone()
		}
	}

	if input.MailingAddress != nil {
		if *input.MailingAddress != "" {
			update = update.SetMailingAddress(*input.MailingAddress)
		} else {
			update = update.ClearMailingAddress()
		}
	}

	if input.FRN != nil {
		if *input.FRN != "" {
			update = update.SetFrn(*input.FRN)
		} else {
			update = update.ClearFrn()
		}
	}

	if input.LicenseClass != nil {
		if *input.LicenseClass != "" {
			update = update.SetLicenseClass(*input.LicenseClass)
		} else {
			update = update.ClearLicenseClass()
		}
	}

	if input.MemberType != nil {
		if *input.MemberType != "" {
			update = update.SetMemberType(*input.MemberType)
		} else {
			update = update.ClearMemberType()
		}
	}

	if input.IsActive != nil {
		update = update.SetIsActive(*input.IsActive)
	}

	if input.IsSilentKey != nil {
		update = update.SetIsSilentKey(*input.IsSilentKey)
	}

	return update.Save(ctx)
}
