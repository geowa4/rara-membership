package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/geowa4/rara-membership/services"
	"github.com/spf13/cobra"
)

var createMemberCmd = &cobra.Command{
	Use:   "create-member",
	Short: "Create a new member",
	Long:  "Create a new member using an interactive form.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			name     string
			email    string
			phone    string
			callSign string
			frn      string
			isActive bool = true
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
					Validate(func(str string) error {
						if str == "" {
							return fmt.Errorf("email is required")
						}
						// Basic email validation
						if len(str) < 3 || !contains(str, "@") {
							return fmt.Errorf("please enter a valid email address")
						}
						return nil
					}),

				huh.NewInput().
					Title("Phone Number").
					Description("Enter the member's phone number").
					Value(&phone).
					Placeholder("(555) 123-4567"),

				huh.NewInput().
					Title("Call Sign").
					Description("Enter the member's amateur radio call sign").
					Value(&callSign).
					Placeholder("W1ABC"),

				huh.NewInput().
					Title("FRN").
					Description("Enter the member's FCC Registration Number").
					Value(&frn).
					Placeholder("0012345678"),

				huh.NewConfirm().
					Title("Active Member?").
					Description("Is this member currently active?").
					Value(&isActive),
			),
		)

		err := form.Run()
		if err != nil {
			return fmt.Errorf("form error: %w", err)
		}

		// Create the member
		member, err := services.CreateMember(name, email, phone, callSign, frn, isActive)
		if err != nil {
			return fmt.Errorf("failed to create member: %w", err)
		}

		fmt.Printf("\n✅ Member created successfully!\n")
		fmt.Printf("ID: %d\n", member.ID)
		fmt.Printf("Name: %s\n", member.Name)
		fmt.Printf("Email: %s\n", member.Email)
		fmt.Printf("Call Sign: %s\n", member.CallSign)
		fmt.Printf("Active: %t\n", member.IsActive)

		return nil
	},
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}