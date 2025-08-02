package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rara-membership",
	Short: "A CLI tool for managing RARA membership",
	Long:  "A command line interface for managing RARA membership database operations.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listMembersCmd)
	rootCmd.AddCommand(createMemberCmd)
	rootCmd.AddCommand(serveCmd)
}
