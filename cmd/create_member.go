package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/geowa4/rara-membership/services"
	"github.com/geowa4/rara-membership/validation"
	"github.com/spf13/cobra"
)

var createMemberCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new member",
	Long:  "Create a new member using an interactive form.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			name           string
			email          string
			phone          string
			mailingAddress string
			callSign       string
			frn            string
			isActive       bool = true
			isSilentKey    bool = false
			licenseClass   string
			memberType     string
		)

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
					Description("Enter the member's mailing address (optional)").
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

				huh.NewSelect[string]().
					Title("Member Type").
					Description("Select the member type (optional)").
					Options(
						huh.NewOption("None", ""),
						huh.NewOption("Student", "Student"),
						huh.NewOption("Regular", "Regular"),
						huh.NewOption("Senior", "Senior"),
						huh.NewOption("Associate", "Associate"),
					).
					Value(&memberType),

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

		err := form.Run()
		if err != nil {
			return fmt.Errorf("form error: %w", err)
		}

		// Create the member
		member, err := services.CreateMember(name, email, phone, mailingAddress, callSign, frn, isActive, isSilentKey, licenseClass, memberType)
		if err != nil {
			return fmt.Errorf("failed to create member: %w", err)
		}

		fmt.Printf("\n✅ Member created successfully!\n")
		fmt.Printf("ID: %d\n", member.ID)
		fmt.Printf("Name: %s\n", member.Name)
		fmt.Printf("Email: %s\n", member.Email)
		fmt.Printf("Phone: %s\n", member.Phone)
		if member.MailingAddress != "" {
			fmt.Printf("Mailing Address: %s\n", member.MailingAddress)
		}
		if member.CallSign != "" {
			fmt.Printf("Call Sign: %s\n", member.CallSign)
		}
		fmt.Printf("FRN: %s\n", member.Frn)
		if member.LicenseClass != "" {
			fmt.Printf("License Class: %s\n", member.LicenseClass)
		}
		if member.MemberType != "" {
			fmt.Printf("Member Type: %s\n", member.MemberType)
		}
		fmt.Printf("Active: %t\n", member.IsActive)
		fmt.Printf("Silent Key: %t\n", member.IsSilentKey)

		return nil
	},
}
