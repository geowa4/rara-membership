package cmd

import (
	"github.com/geowa4/rara-membership/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  "Start an HTTP server with API endpoints for member management.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Start()
	},
}
