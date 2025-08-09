package cmd

import (
	"github.com/spf13/cobra"
)

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Manage events",
	Long:  "Commands for managing events in the RARA membership system.",
}

func init() {
	eventCmd.AddCommand(listEventsCmd)
	eventCmd.AddCommand(createEventCmd)
	eventCmd.AddCommand(updateEventCmd)
}
