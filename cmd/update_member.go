package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/geowa4/rara-membership/services"
	"github.com/geowa4/rara-membership/validation"
	"github.com/spf13/cobra"
)

var updateMemberCmd = &cobra.Command{
	Use:   "update [call_sign]",
	Short: "Update an existing member",
	Long:  "Update an existing active member using their call sign.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		callSignArg := args[0]

		// Find the member by call sign
		member, err := services.GetActiveMemberByCallSign(callSignArg)
		if err != nil {
			return fmt.Errorf("failed to find active member with call sign %s: %w", callSignArg, err)
		}

		// Pre-populate form with current values
		var (
			name           = member.Name
			email          = member.Email
			phone          = member.Phone
			mailingAddress = member.MailingAddress
			callSign       = member.CallSign
			frn            = member.Frn
			isActive       = member.IsActive
			isSilentKey    = member.IsSilentKey
			licenseClass   = member.LicenseClass
		)

		fmt.Printf("📝 Updating member: %s (%s)\n\n", member.Name, member.CallSign)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Full Name").
					Description("Enter the member's full name").
					Value(&name).
					Validate(func(str string) error {
						if str == "" {
							return fmt.Errorf("name is required")
						}
						return nil
					}),

				huh.NewInput().
					Title("Email Address").
					Description("Enter the member's email address").
					Value(&email).
					Validate(validation.ValidateEmail),

				huh.NewInput().
					Title("Phone Number").
					Description("Enter the member's phone number").
					Value(&phone).
					Placeholder("(555) 123-4567"),

				huh.NewInput().
					Title("Mailing Address").
					Description("Enter the member's mailing address").
					Value(&mailingAddress).
					Placeholder("123 Main St, City, State 12345"),

				huh.NewInput().
					Title("Call Sign").
					Description("Enter the member's amateur radio call sign (optional)").
					Value(&callSign).
					Placeholder("N0CALL"),

				huh.NewInput().
					Title("FRN").
					Description("Enter the member's FCC Registration Number").
					Value(&frn).
					Validate(func(str string) error {
						if str == "" {
							return fmt.Errorf("FRN is required")
						}
						return nil
					}).
					Placeholder("0012345678"),

				huh.NewSelect[string]().
					Title("License Class").
					Description("Select the member's amateur radio license class (optional)").
					Options(
						huh.NewOption("None", ""),
						huh.NewOption("Technician", "Technician"),
						huh.NewOption("General", "General"),
						huh.NewOption("Extra", "Extra"),
						huh.NewOption("Novice", "Novice"),
						huh.NewOption("Advanced", "Advanced"),
					).
					Value(&licenseClass),

				huh.NewConfirm().
					Title("Active Member?").
					Description("Is this member currently active?").
					Value(&isActive),

				huh.NewConfirm().
					Title("Silent Key?").
					Description("Is this member a silent key (deceased)?").
					Value(&isSilentKey),
			),
		)

		err = form.Run()
		if err != nil {
			return fmt.Errorf("form error: %w", err)
		}

		// Update the member
		updatedMember, err := services.UpdateMember(member.ID, name, email, phone, mailingAddress, callSign, frn, isActive, isSilentKey, licenseClass)
		if err != nil {
			return fmt.Errorf("failed to update member: %w", err)
		}

		fmt.Printf("\n✅ Member updated successfully!\n")
		fmt.Printf("ID: %d\n", updatedMember.ID)
		fmt.Printf("Name: %s\n", updatedMember.Name)
		fmt.Printf("Email: %s\n", updatedMember.Email)
		fmt.Printf("Phone: %s\n", updatedMember.Phone)
		fmt.Printf("Address: %s\n", updatedMember.MailingAddress)
		if updatedMember.CallSign != "" {
			fmt.Printf("Call Sign: %s\n", updatedMember.CallSign)
		}
		fmt.Printf("FRN: %s\n", updatedMember.Frn)
		if updatedMember.LicenseClass != "" {
			fmt.Printf("License Class: %s\n", updatedMember.LicenseClass)
		}
		fmt.Printf("Active: %t\n", updatedMember.IsActive)
		fmt.Printf("Silent Key: %t\n", updatedMember.IsSilentKey)

		return nil
	},
}
