package cmd

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/services"
	"github.com/spf13/cobra"
)

// PointTransaction represents either an allocation or deduction for display
type PointTransaction struct {
	Type        string    // "EARNED" or "REDEEMED"
	Points      int       // positive for earned, negative for redeemed (for balance calculation)
	DisplayDate time.Time // when the transaction occurred
	Description string    // event name or redemption description
	Notes       string    // additional notes
}

var pointsOnly bool

// memberPointsCmd represents the member points command
var memberPointsCmd = &cobra.Command{
	Use:   "points [call_sign]",
	Short: "Show complete point history for a member by call sign",
	Long: `Display complete point history for a member, including both points earned 
from volunteering at events and points redeemed for rewards. Shows running balance
and chronological transaction history. Only shows history for active members.

Examples:
  # Show complete point history for member with call sign W1ABC
  rara-membership member points W1ABC
  
  # Show complete point history for member with call sign KD2DEF
  rara-membership member points KD2DEF
  
  # Show only the current balance as a plain number
  rara-membership member points W1ABC --points-only`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		callSign := strings.ToUpper(strings.TrimSpace(args[0]))

		if callSign == "" {
			return fmt.Errorf("call sign cannot be empty")
		}

		// Find the member by call sign using service function
		mbr, err := services.GetActiveMemberByCallSign(callSign)
		if err != nil {
			return fmt.Errorf("member with call sign %s not found or is not active", callSign)
		}

		// Create services
		allocationService := services.NewPointAllocationService(database.Client)
		deductionService := services.NewPointDeductionService(database.Client)

		// Get allocations (points earned)
		allocations, err := allocationService.GetAllocationsByMember(ctx, mbr.ID)
		if err != nil {
			return fmt.Errorf("error getting point allocations: %w", err)
		}

		// Get deductions (points redeemed)
		deductions, err := deductionService.GetDeductionsByMember(ctx, mbr.ID)
		if err != nil {
			return fmt.Errorf("error getting point deductions: %w", err)
		}

		// Convert to unified transaction format
		var transactions []PointTransaction

		// Add allocations
		for _, allocation := range allocations {
			eventName := allocation.Edges.Event.Name
			transactions = append(transactions, PointTransaction{
				Type:        "EARNED",
				Points:      allocation.Points,
				DisplayDate: allocation.CreatedAt,
				Description: eventName,
				Notes:       allocation.Notes,
			})
		}

		// Add deductions
		for _, deduction := range deductions {
			transactions = append(transactions, PointTransaction{
				Type:        "REDEEMED",
				Points:      -deduction.Points, // negative for balance calculation
				DisplayDate: deduction.CreatedAt,
				Description: "Points Redeemed",
				Notes:       deduction.Notes,
			})
		}

		// Calculate current balance
		currentBalance := 0
		for _, transaction := range transactions {
			currentBalance += transaction.Points
		}

		// If points-only output is requested, just print the balance and return
		if pointsOnly {
			fmt.Println(currentBalance)
			return nil
		}

		if len(transactions) == 0 {
			fmt.Printf("No point history found for %s (%s)\n", mbr.Name, callSign)
			return nil
		}

		// Sort transactions by date (oldest first) for chronological display
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].DisplayDate.Before(transactions[j].DisplayDate)
		})

		// Define colors and styles
		var (
			purple    = lipgloss.Color("99")
			gray      = lipgloss.Color("245")
			lightGray = lipgloss.Color("241")
			green     = lipgloss.Color("42")
			red       = lipgloss.Color("196")
			blue      = lipgloss.Color("33")

			headerStyle   = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
			cellStyle     = lipgloss.NewStyle().Padding(0, 1)
			oddRowStyle   = cellStyle.Foreground(gray)
			evenRowStyle  = cellStyle.Foreground(lightGray)
			earnedStyle   = cellStyle.Foreground(green)
			redeemedStyle = cellStyle.Foreground(red)
		)

		// Prepare table data with running balance
		rows := make([][]string, len(transactions))
		runningBalance := 0
		totalEarned := 0
		totalRedeemed := 0

		for i, transaction := range transactions {
			runningBalance += transaction.Points

			if transaction.Points > 0 {
				totalEarned += transaction.Points
			} else {
				totalRedeemed += -transaction.Points
			}

			pointsDisplay := strconv.Itoa(transaction.Points)
			if transaction.Points > 0 {
				pointsDisplay = "+" + pointsDisplay
			}

			notes := transaction.Notes
			if notes == "" {
				notes = "-"
			}

			rows[i] = []string{
				transaction.DisplayDate.Format("2006-01-02"),
				transaction.Type,
				transaction.Description,
				pointsDisplay,
				strconv.Itoa(runningBalance),
				notes,
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
					// Color points column based on transaction type
					if col == 3 && row > 0 { // Points column
						if strings.HasPrefix(rows[row-1][3], "+") {
							return earnedStyle
						} else {
							return redeemedStyle
						}
					}
					return evenRowStyle
				default:
					// Color points column based on transaction type
					if col == 3 && row > 0 { // Points column
						if strings.HasPrefix(rows[row-1][3], "+") {
							return earnedStyle
						} else {
							return redeemedStyle
						}
					}
					return oddRowStyle
				}
			}).
			Headers("DATE", "TYPE", "DESCRIPTION", "POINTS", "BALANCE", "NOTES").
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

		fmt.Println(titleStyle.Render(fmt.Sprintf("🏆 Complete Point History for %s (%s)", mbr.Name, callSign)))
		fmt.Println(t)
		fmt.Println(summaryStyle.Render(fmt.Sprintf("📊 Summary: %d Earned • %d Redeemed • %d Current Balance", totalEarned, totalRedeemed, runningBalance)))

		return nil
	},
}

func init() {
	memberCmd.AddCommand(memberPointsCmd)
	memberPointsCmd.Flags().BoolVar(&pointsOnly, "points-only", false, "Output only the current balance as a plain number")
}
