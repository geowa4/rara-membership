package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/ent/member"
	"github.com/geowa4/rara-membership/services"
	"github.com/spf13/cobra"
)

// memberPointsCmd represents the member points command
var memberPointsCmd = &cobra.Command{
	Use:   "points [call_sign]",
	Short: "Show point history for a member by call sign",
	Long: `Display all events a member has volunteered for, including event name, date, 
points allocated, and when the points were allocated. Only shows history for 
active members (not silent keys).

Examples:
  # Show point history for member with call sign W1ABC
  rara-membership member points W1ABC
  
  # Show point history for member with call sign KD2DEF
  rara-membership member points KD2DEF`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		callSign := strings.ToUpper(strings.TrimSpace(args[0]))
		
		if callSign == "" {
			return fmt.Errorf("call sign cannot be empty")
		}

		// Find the member by call sign (case-insensitive search)
		mbr, err := database.Client.Member.Query().
			Where(
				member.Or(member.CallSignEQ(strings.ToUpper(callSign)), member.CallSignEQ(strings.ToLower(callSign))),
				member.IsSilentKeyEQ(false),
			).
			Only(ctx)
		if err != nil {
			return fmt.Errorf("member with call sign %s not found or is a silent key", callSign)
		}

		// Create service using the global database client
		pointService := services.NewPointAllocationService(database.Client)

		// Get volunteer history
		allocations, err := pointService.GetVolunteerHistoryForMember(ctx, mbr.ID)
		if err != nil {
			return fmt.Errorf("error getting volunteer history: %w", err)
		}

		if len(allocations) == 0 {
			fmt.Printf("No volunteer history found for %s (%s)\n", mbr.Name, callSign)
			return nil
		}

		// Define colors and styles
		var (
			purple    = lipgloss.Color("99")
			gray      = lipgloss.Color("245")
			lightGray = lipgloss.Color("241")
			green     = lipgloss.Color("42")
			blue      = lipgloss.Color("33")

			headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
			cellStyle    = lipgloss.NewStyle().Padding(0, 1)
			oddRowStyle  = cellStyle.Foreground(gray)
			evenRowStyle = cellStyle.Foreground(lightGray)
		)

		// Prepare table data
		rows := make([][]string, len(allocations))
		totalPoints := 0

		for i, allocation := range allocations {
			eventName := allocation.Edges.Event.Name
			eventDate := allocation.Edges.Event.Date.Format("2006-01-02")
			points := allocation.Points
			allocatedOn := allocation.CreatedAt.Format("2006-01-02")
			notes := allocation.Notes
			
			if notes == "" {
				notes = "-"
			}

			rows[i] = []string{
				eventName,
				eventDate,
				strconv.Itoa(points),
				allocatedOn,
				notes,
			}
			
			totalPoints += points
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
			Headers("EVENT NAME", "EVENT DATE", "POINTS", "ALLOCATED ON", "NOTES").
			Rows(rows...)

		// Print title and table
		titleStyle := lipgloss.NewStyle().
			Foreground(green).
			Bold(true).
			MarginBottom(1)

		summaryStyle := lipgloss.NewStyle().
			Foreground(blue).
			Bold(true).
			MarginTop(1)

		fmt.Println(titleStyle.Render(fmt.Sprintf("🏆 Volunteer History for %s (%s)", mbr.Name, callSign)))
		fmt.Println(t)
		fmt.Println(summaryStyle.Render(fmt.Sprintf("📊 Summary: %d Events • %d Total Points", len(allocations), totalPoints)))

		return nil
	},
}


func init() {
	memberCmd.AddCommand(memberPointsCmd)
}