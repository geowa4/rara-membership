package cmd

import (
	"context"
	"fmt"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/services"
	"github.com/spf13/cobra"
)

var (
	eventID  int
	memberID int
	points   int
	useDefaultPoints bool
	notes    string
)

// volunteerCmd represents the volunteer command
var volunteerCmd = &cobra.Command{
	Use:   "volunteer",
	Short: "Allocate points to a member for volunteering at an event",
	Long: `Allocate points to a member for volunteering at an event.
	
If no points are specified, the default points for the event will be used.

Examples:
  # Allocate default points to member 1 for event 5
  rara-membership volunteer --event 5 --member 1
  
  # Allocate specific points with a note
  rara-membership volunteer --event 5 --member 1 --points 20 --notes "Led the event setup"
  
  # Use flags shorthand
  rara-membership volunteer -e 5 -m 1 -p 15`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		
		// Create service using the global database client
		pointService := services.NewPointAllocationService(database.Client)

		// Determine points to allocate
		var pointsPtr *int
		if !useDefaultPoints {
			pointsPtr = &points
		}

		// Allocate points
		allocation, err := pointService.AllocatePoints(ctx, eventID, memberID, pointsPtr, notes)
		if err != nil {
			return fmt.Errorf("error allocating points: %w", err)
		}

		// Display success message
		fmt.Printf("\n✅ Points allocated successfully!\n")
		fmt.Printf("Member: %s\n", allocation.Edges.Member.Name)
		fmt.Printf("Event: %s\n", allocation.Edges.Event.Name)
		fmt.Printf("Points: %d\n", allocation.Points)
		
		if allocation.Notes != "" {
			fmt.Printf("Notes: %s\n", allocation.Notes)
		}

		// Get and display total points for the member
		totalPoints, err := pointService.GetTotalPointsForMember(ctx, memberID)
		if err != nil {
			fmt.Printf("Warning: Could not get total points: %v\n", err)
		} else {
			fmt.Printf("Total points earned by member: %d\n", totalPoints)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(volunteerCmd)

	volunteerCmd.Flags().IntVarP(&eventID, "event", "e", 0, "Event ID (required)")
	volunteerCmd.Flags().IntVarP(&memberID, "member", "m", 0, "Member ID (required)")
	volunteerCmd.Flags().IntVarP(&points, "points", "p", 0, "Points to allocate (optional, uses event default if not specified)")
	volunteerCmd.Flags().StringVarP(&notes, "notes", "n", "", "Optional notes about the allocation")

	volunteerCmd.MarkFlagRequired("event")
	volunteerCmd.MarkFlagRequired("member")

	// Custom flag parsing to handle optional points
	volunteerCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Check if points flag was explicitly set
		if cmd.Flags().Changed("points") {
			useDefaultPoints = false
			if points < 0 || points > 100 {
				return fmt.Errorf("points must be between 0 and 100")
			}
		} else {
			useDefaultPoints = true
		}

		// Validate IDs
		if eventID <= 0 {
			return fmt.Errorf("event ID must be a positive number")
		}
		if memberID <= 0 {
			return fmt.Errorf("member ID must be a positive number")
		}

		return nil
	}
}