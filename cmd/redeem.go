package cmd

import (
	"context"
	"fmt"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/services"
	"github.com/spf13/cobra"
)

var (
	redeemMemberID int
	redeemPoints   int
	redeemNotes    string
)

// redeemCmd represents the redeem command
var redeemCmd = &cobra.Command{
	Use:   "redeem",
	Short: "Redeem points for a member",
	Long: `Redeem points for a member, creating a point deduction record.
	
This command allows you to deduct points from a member's balance when they
redeem their points for rewards or benefits. The system ensures that the
member has sufficient points before allowing the redemption.

Examples:
  # Redeem 25 points for member 1
  rara-membership redeem --member 1 --points 25
  
  # Redeem points with a note about what was redeemed
  rara-membership redeem --member 1 --points 50 --notes "T-shirt reward"
  
  # Use flags shorthand
  rara-membership redeem -m 1 -p 30 -n "Club merchandise"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Create service using the global database client
		deductionService := services.NewPointDeductionService(database.Client)

		// Check current balance before attempting redemption
		balance, err := deductionService.GetPointBalanceForMember(ctx, redeemMemberID)
		if err != nil {
			return fmt.Errorf("error checking member point balance: %w", err)
		}

		fmt.Printf("Current point balance: %d\n", balance)

		if balance < redeemPoints {
			return fmt.Errorf("insufficient points: member has %d points but trying to redeem %d", balance, redeemPoints)
		}

		// Redeem points
		deduction, err := deductionService.DeductPoints(ctx, redeemMemberID, redeemPoints, redeemNotes)
		if err != nil {
			return fmt.Errorf("error redeeming points: %w", err)
		}

		// Display success message
		fmt.Printf("\n✅ Points redeemed successfully!\n")
		fmt.Printf("Member: %s\n", deduction.Edges.Member.Name)
		fmt.Printf("Points redeemed: %d\n", deduction.Points)

		if deduction.Notes != "" {
			fmt.Printf("Notes: %s\n", deduction.Notes)
		}

		// Get and display new balance
		newBalance, err := deductionService.GetPointBalanceForMember(ctx, redeemMemberID)
		if err != nil {
			fmt.Printf("Warning: Could not get updated balance: %v\n", err)
		} else {
			fmt.Printf("Remaining point balance: %d\n", newBalance)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(redeemCmd)

	redeemCmd.Flags().IntVarP(&redeemMemberID, "member", "m", 0, "Member ID (required)")
	redeemCmd.Flags().IntVarP(&redeemPoints, "points", "p", 0, "Points to redeem (required)")
	redeemCmd.Flags().StringVarP(&redeemNotes, "notes", "n", "", "Optional notes about the redemption")

	redeemCmd.MarkFlagRequired("member")
	redeemCmd.MarkFlagRequired("points")

	// Validation
	redeemCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Validate member ID
		if redeemMemberID <= 0 {
			return fmt.Errorf("member ID must be a positive number")
		}

		// Validate points
		if redeemPoints <= 0 {
			return fmt.Errorf("points must be a positive number")
		}

		return nil
	}
}
