package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

		// Define colors and styles
		var (
			purple    = lipgloss.Color("99")
			gray      = lipgloss.Color("245")
			lightGray = lipgloss.Color("241")
			green     = lipgloss.Color("42")

			headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
			cellStyle    = lipgloss.NewStyle().Padding(0, 1)
			oddRowStyle  = cellStyle.Foreground(gray)
			evenRowStyle = cellStyle.Foreground(lightGray)
		)

		// Prepare member data for table
		activeStatus := "✓"
		if !updatedMember.IsActive {
			activeStatus = "✗"
		}
		silentKeyStatus := ""
		if updatedMember.IsSilentKey {
			silentKeyStatus = "sk"
		}
		memberCallSign := updatedMember.CallSign
		if memberCallSign == "" {
			memberCallSign = "-"
		}
		memberLicenseClass := updatedMember.LicenseClass
		if memberLicenseClass == "" {
			memberLicenseClass = "-"
		}

		row := []string{
			strconv.Itoa(updatedMember.ID),
			updatedMember.Name,
			memberCallSign,
			updatedMember.Email,
			updatedMember.Phone,
			updatedMember.Frn,
			memberLicenseClass,
			activeStatus,
			silentKeyStatus,
		}

		// Create and style the table
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(purple)).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == table.HeaderRow:
					return headerStyle
				case row%2 == 0:
					return evenRowStyle
				default:
					return oddRowStyle
				}
			}).
			Headers("ID", "NAME", "CALL SIGN", "EMAIL", "PHONE", "FRN", "LICENSE", "ACTIVE", "SILENT KEY").
			Rows([][]string{row}...)

		// Print title and table
		titleStyle := lipgloss.NewStyle().
			Foreground(green).
			Bold(true).
			MarginBottom(1)

		fmt.Println()
		fmt.Println(titleStyle.Render("✅ Member Updated Successfully!"))
		fmt.Println(t)

		return nil
	},
}
