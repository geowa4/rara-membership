package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/geowa4/rara-membership/services"
	"github.com/spf13/cobra"
)

var listMembersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active members",
	Long:  "Display a list of all active members in the database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		members, err := services.GetActiveMembers()
		if err != nil {
			return fmt.Errorf("failed querying members: %w", err)
		}

		if len(members) == 0 {
			fmt.Println("No active members found.")
			return nil
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

		// Prepare table data
		rows := make([][]string, len(members))
		for i, m := range members {
			activeStatus := "✓"
			if !m.IsActive {
				activeStatus = "✗"
			}
			silentKeyStatus := ""
			if m.IsSilentKey {
				silentKeyStatus = "sk"
			}
			callSign := m.CallSign
			if callSign == "" {
				callSign = "-"
			}
			licenseClass := m.LicenseClass
			if licenseClass == "" {
				licenseClass = "-"
			}
			rows[i] = []string{
				strconv.Itoa(m.ID),
				m.Name,
				callSign,
				m.Email,
				m.Phone,
				m.Frn,
				licenseClass,
				activeStatus,
				silentKeyStatus,
			}
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
			Rows(rows...)

		// Print title and table
		titleStyle := lipgloss.NewStyle().
			Foreground(green).
			Bold(true).
			MarginBottom(1)

		fmt.Println(titleStyle.Render(fmt.Sprintf("📻 Active Members (%d)", len(members))))
		fmt.Println(t)

		return nil
	},
}
