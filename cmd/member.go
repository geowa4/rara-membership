package cmd

import (
	"github.com/spf13/cobra"
)

var memberCmd = &cobra.Command{
	Use:   "member",
	Short: "Member management commands",
	Long:  "Commands for managing RARA members including listing, creating, and updating members.",
}

func init() {
	memberCmd.AddCommand(listMembersCmd)
	memberCmd.AddCommand(createMemberCmd)
	memberCmd.AddCommand(updateMemberCmd)
}
